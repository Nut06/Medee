# -----------------------------------------------
# AWS Load Balancer Controller
# watches Gateway/HTTPRoute objects → creates ALB
# -----------------------------------------------
resource "helm_release" "aws_load_balancer_controller" {
  count = var.enable_gateway_api ? 1 : 0

  name       = "aws-load-balancer-controller"
  repository = "https://aws.github.io/eks-charts"
  chart      = "aws-load-balancer-controller"
  version    = "3.1.0"
  namespace  = "kube-system"

  wait            = true
  wait_for_jobs   = true
  timeout         = 600
  atomic          = true
  cleanup_on_fail = true

  values = [yamlencode({
    clusterName = var.cluster_name
    region      = var.aws_region
    vpcId = var.vpc_id

    serviceAccount = {
      create = true
      name   = "aws-load-balancer-controller"
      annotations = {
        "eks.amazonaws.com/role-arn" = var.lb_controller_role_arn
      }
    }

    enableGatewayAPI = true

    resources = {
      limits   = { cpu = "200m", memory = "500Mi" }
      requests = { cpu = "100m", memory = "200Mi" }
    }

    nodeSelector        = { "kubernetes.io/os" = "linux" }
    replicaCount        = 2
    podDisruptionBudget = { maxUnavailable = 1 }
  })]

  # depends_on = [helm_release.gateway_api[0]]
}

resource "null_resource" "gatewayAPI_crds" {
  count = var.enable_gateway_api ? 1 : 0
  provisioner "local-exec" {
    command = "kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.3.0/standard-install.yaml"
  }

  # depends_on = [helm_release.aws_load_balancer_controller[0]]
}

resource "null_resource" "lbc_gateway_crds" {
  count = var.enable_gateway_api ? 1 : 0
  provisioner "local-exec" {
    command = "kubectl apply -f https://raw.githubusercontent.com/kubernetes-sigs/aws-load-balancer-controller/refs/heads/main/config/crd/gateway/gateway-crds.yaml"
  }

  depends_on = [helm_release.aws_load_balancer_controller[0]]
}

resource "helm_release" "argocd" {
  name       = "argocd"
  repository = "https://argoproj.github.io/argo-helm"
  chart      = "argo-cd"
  version    = "9.4.16"
  namespace  = "argocd"

  create_namespace = true
  wait          = true
  wait_for_jobs = true
  timeout       = 600
  atomic        = true
  cleanup_on_fail = true

  values = [
    yamlencode({
      server = {
        service = {
          type = "ClusterIP"
        }
        ingress = {
          enabled = false
        }
      }
      # Minimal configuration - detailed config goes to GitOps
    })
  ]
}
# -----------------------------------------------
# gateway-system namespace
# -----------------------------------------------
resource "kubernetes_namespace_v1" "gateway_system" {
  count = var.enable_gateway_api ? 1 : 0

  metadata {
    name = "gateway-system"
    labels = {
      name                           = "gateway-system"
      "app.kubernetes.io/managed-by" = "terraform"
      environment                    = var.environment
      project                        = var.project_name
    }
  }

  # depends_on = [helm_release.aws_load_balancer_controller[0]]
}


# start comment here 2nd step in this file
# -----------------------------------------------
# GatewayClass
# registers ALB as the controller for Gateway objects
# -----------------------------------------------
resource "kubernetes_manifest" "gateway_class" {
  count = var.enable_gateway_api ? 1 : 0

  manifest = {
    apiVersion = "gateway.networking.k8s.io/v1beta1"
    kind       = "GatewayClass"
    metadata = {
      name   = var.gateway_class_name
      labels = { "app.kubernetes.io/managed-by" = "terraform" }
    }
    spec = {
      controllerName = "gateway.k8s.aws/alb"
      description    = "AWS Load Balancer Controller Gateway Class"
    }
  }

  # depends_on = [
  #   helm_release.aws_load_balancer_controller[0],
  #   helm_release.gateway_api[0],
  # ]
}


# -----------------------------------------------
# Main Gateway
# single ALB entry point for all traffic to dodee.me
# handles: frontend, backend/api, argocd, grafana
# -----------------------------------------------
resource "kubernetes_manifest" "main_gateway" {
  count = var.enable_gateway_api ? 1 : 0

  manifest = {
    apiVersion = "gateway.networking.k8s.io/v1beta1"
    kind       = "Gateway"
    metadata = {
      name      = "main-gateway"
      namespace = "gateway-system"
      labels    = { "app.kubernetes.io/managed-by" = "terraform" }
    }
    spec = {
      gatewayClassName = var.gateway_class_name
      infrastructure = {
        parametersRef = {
          kind  = "LoadBalancerConfiguration"
          name  = "main-gateway-lbconfig"
          group = "gateway.k8s.aws"
        }
      }
      listeners = [
        {
          name          = "http"
          protocol      = "HTTP"
          port          = 80
          allowedRoutes = { namespaces = { from = "All" } }
        },
        {
          name          = "https"
          protocol      = "HTTPS"
          port          = 443
          allowedRoutes = { namespaces = { from = "All" } }
        }
      ]
    }
  }

  depends_on = [
    kubernetes_manifest.gateway_class[0],
    kubernetes_namespace_v1.gateway_system[0],
    kubernetes_manifest.main_gateway_lbconfig[0],
  ]
}

resource "kubernetes_manifest" "main_gateway_lbconfig" {
  count = var.enable_gateway_api ? 1 : 0

  manifest = {
    apiVersion = "gateway.k8s.aws/v1beta1"
    kind       = "LoadBalancerConfiguration"
    metadata = {
      name      = "main-gateway-lbconfig"
      namespace = "gateway-system"
    }
    spec = {
      scheme = "internet-facing"
      loadBalancerName = "${var.project_name}-main-gateway"
      loadBalancerSubnets = [
        { identifier = var.public_subnet_ids[0] },
        { identifier = var.public_subnet_ids[1] },
      ]

      listenerConfigurations = [
        {
          protocolPort       = "HTTPS:443"
          defaultCertificate = var.acm_certificate_arn
          sslPolicy = "ELBSecurityPolicy-TLS13-1-2-2021-06"
        }
      ]

      tags = {
        Environment = var.environment
        Project     = var.project_name
        ManagedBy   = "terraform"
      }
    }
  }

  depends_on = [
    null_resource.lbc_gateway_crds[0],
    kubernetes_namespace_v1.gateway_system[0]
  ]
}


resource "helm_release" "prometheus" {
  name       = "prometheus"
  repository = "https://prometheus-community.github.io/helm-charts"
  chart      = "kube-prometheus-stack"
  version    = "82.18.0"
  namespace  = "monitoring"
  
  create_namespace = true
  
  values = [yamlencode({
    prometheus = {
      prometheusSpec = {
        storageSpec = {
          volumeClaimTemplate = {
            spec = {
              storageClassName = "gp3"
              resources = {
                requests = {
                  storage = "10Gi"
                }
              }
            }
          }
        }
      }
    }
    
    grafana = {
      enabled = true
      adminPassword = "admin"
      persistence = {
        enabled = true
        storageClassName = "gp3"
      }
    }
    
    alertmanager = {
      alertmanagerSpec = {
        storage = {
          volumeClaimTemplate = {
            spec = {
              storageClassName = "gp3"
            }
          }
        }
      }
  }
})]
}

resource "helm_release" "vault" {
  name       = "vault"
  repository = "https://helm.releases.hashicorp.com"
  chart      = "vault"
  namespace  = "vault"

  create_namespace = true

  values = [yamlencode({
    server = {
      ha = {
        enabled = false
      }

      dataStorage = {
        enabled = true
        size = "5Gi"
        storageClass = "gp3"  # ← เพิ่มบรรทัดนี้
      }
      
    }
    injector = {
      enabled = true
    }
  })]
}
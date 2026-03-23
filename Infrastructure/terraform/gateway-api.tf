resource "kubernetes_manifest" "gateway_api_crds" {
  for_each = toset([
    "gatewayclasses",
    "gateways", 
    "httproutes",
    "referencegrants",
    "tcproutes",
    "tlsroutes",
    "udproutes"
  ])

  manifest = {
    apiVersion = "apiextensions.k8s.io/v1"
    kind       = "CustomResourceDefinition"
    metadata = {
      name = "${each.value}.gateway.networking.k8s.io"
    }
    spec = {
      group = "gateway.networking.k8s.io"
      versions = [
        {
          name    = "v1"
          served  = true
          storage = true
          schema = {
            openAPIV3Schema = {
              type = "object"
              properties = {
                spec = {
                  type = "object"
                }
                status = {
                  type = "object"
                }
              }
            }
          }
        }
      ]
      scope = "Namespaced"
      names = {
        plural   = each.value
        singular = replace(each.value, "s$", "")
        kind     = title(replace(each.value, "s$", ""))
      }
    }
  }

  depends_on = [aws_eks_cluster.medee-cluster]
}

resource "helm_release" "aws_load_balancer_controller" {
  name       = "aws-load-balancer-controller"
  repository = "https://aws.github.io/eks-charts"
  chart      = "aws-load-balancer-controller"
  version    = "1.8.1"
  namespace  = "kube-system"

  values = [
    yamlencode({
      clusterName = var.cluster_name
      region = var.aws_region

      serviceAccount = {
        create = true
        name   = "aws-load-balancer-controller"
        annotations = {
          "eks.amazonaws.com/role-arn" = aws_iam_role.aws_load_balancer_controller[0].arn
        }
      }

      # Enable Gateway API support
      enableGatewayAPI = true
      
      # Resource configuration
      resources = {
        limits = {
          cpu    = "200m"
          memory = "500Mi"
        }
        requests = {
          cpu    = "100m"
          memory = "200Mi"
        }
      }

      # Node selector for running on EC2 nodes
      nodeSelector = {
        "kubernetes.io/os" = "linux"
      }

      # Tolerations if needed
      tolerations = []
    })
  ]

  depends_on = [
    aws_iam_role_policy_attachment.aws_load_balancer_controller[0],
    kubernetes_manifest.gateway_api_crds
  ]
}

resource "helm_release" "cert_manager" {
  name = "cert-manager"
  repository = "https://charts.jetstack.io"
  chart      = "cert-manager"
  version    = "v1.15.3"
  namespace  = "cert-manager"

  create_namespace = true

  values = [
    yamlencode({
      installCRDs = true
      
      # Resource configuration
      resources = {
        limits = {
          cpu    = "100m"
          memory = "128Mi"
        }
        requests = {
          cpu    = "50m"
          memory = "64Mi"
        }
      }

      webhook = {
        resources = {
          limits = {
            cpu    = "100m"
            memory = "128Mi"
          }
          requests = {
            cpu    = "50m"
            memory = "64Mi"
          }
        }
      }

      cainjector = {
        resources = {
          limits = {
            cpu    = "100m"
            memory = "128Mi"
          }
          requests = {
            cpu    = "50m"
            memory = "64Mi"
          }
        }
      }

      # Node selector
      nodeSelector = {
        "kubernetes.io/os" = "linux"
      }
    })
  ]

  depends_on = [aws_eks_node_group.medee-node-group]
}

resource "kubernetes_manifest" "letsencrypt_cluster_issuer" {
  manifest = {
    apiVersion = "cert-manager.io/v1"
    kind       = "ClusterIssuer"
    metadata = {
      name = "letsencrypt-prod"
    }
    spec = {
      acme = {
        server = "https://acme-v02.api.letsencrypt.org/directory"
        email  = "admin@${var.domain_name}"
        privateKeySecretRef = {
          name = "letsencrypt-prod"
        }
        solvers = [
          {
            http01 = {
              ingress = {
                class = "aws-load-balancer-controller"
              }
            }
          }
        ]
      }
    }
  }

  depends_on = [helm_release.cert_manager[0]]
}

resource "kubernetes_manifest" "gateway_class" {
   manifest = {
    apiVersion = "gateway.networking.k8s.io/v1"
    kind       = "GatewayClass"
    metadata = {
      name = var.gateway_class_name
      labels = {
        "app.kubernetes.io/managed-by" = "terraform"
      }
    }
    spec = {
      controllerName = "ingress.k8s.aws/alb"
      description    = "AWS Load Balancer Controller Gateway Class"
    }
  }

  depends_on = [
    helm_release.aws_load_balancer_controller[0],
    kubernetes_manifest.gateway_api_crds
  ]
}

resource "kubernetes_manifest" "main_gateway" {

  manifest = {
    apiVersion = "gateway.networking.k8s.io/v1"
    kind       = "Gateway"
    metadata = {
      name      = "main-gateway"
      namespace = "gateway-system"
      labels = {
        "app.kubernetes.io/managed-by" = "terraform"
      }
      annotations = {
        "alb.ingress.kubernetes.io/scheme"      = "internet-facing"
        "alb.ingress.kubernetes.io/target-type" = "ip"
        "alb.ingress.kubernetes.io/listen-ports" = jsonencode([
          { HTTP = 80 },
          { HTTPS = 443 }
        ])
      }
    }
    spec = {
      gatewayClassName = var.gateway_class_name
      listeners = [
        {
          name     = "http"
          protocol = "HTTP"
          port     = 80
          allowedRoutes = {
            namespaces = {
              from = "All"
            }
          }
        },
        {
          name     = "https"
          protocol = "HTTPS"
          port     = 443
          tls = {
            mode = "Terminate"
            certificateRefs = var.enable_cert_manager ? [
              {
                name = "main-gateway-tls"
                kind = "Secret"
              }
            ] : []
          }
          allowedRoutes = {
            namespaces = {
              from = "All"
            }
          }
        }
      ]
    }
  }

  depends_on = [
    kubernetes_manifest.gateway_class[0],
    kubernetes_namespace.gateway_system
  ]
}

resource "kubernetes_namespace" "gateway_system" {
  count = var.enable_gateway_api ? 1 : 0

  metadata {
    name = "gateway-system"
    
    labels = {
      name                          = "gateway-system"
      "app.kubernetes.io/managed-by" = "terraform"
      environment                   = var.environment
      project                       = var.project_name
      purpose                       = "gateway"
    }
    
    annotations = {
      "terraform.io/managed" = "true"
    }
  }

  depends_on = [aws_eks_cluster.medee-cluster]
}

resource "kubernetes_manifest" "reference_grants" {
  for_each = var.enable_gateway_api ? toset([
    "argocd",
    var.monitoring_namespace,
    "frontend-${var.environment}",
    "backend-${var.environment}"
  ]) : toset([])

  manifest = {
    apiVersion = "gateway.networking.k8s.io/v1beta1"
    kind       = "ReferenceGrant"
    metadata = {
      name      = "gateway-reference-grant"
      namespace = each.value
    }
    spec = {
      from = [
        {
          group     = "gateway.networking.k8s.io"
          kind      = "HTTPRoute"
          namespace = each.value
        }
      ]
      to = [
        {
          group = "gateway.networking.k8s.io"
          kind  = "Gateway"
        }
      ]
    }
  }

  depends_on = [
    kubernetes_manifest.main_gateway[0],
    kubernetes_namespace.app_namespaces,
    kubernetes_namespace.monitoring,
    kubernetes_namespace.argocd
  ]
}
# -----------------------------------------------
# Gateway API CRDs
# -----------------------------------------------
resource "helm_release" "gateway_api" {
  count = var.enable_gateway_api ? 1 : 0

  name      = "gateway-api"
  chart     = "oci://ghcr.io/nicklasfrahm/charts/gateway-api"
  version   = "0.2.0"
  namespace = "gateway-system"

  create_namespace = true
}

# -----------------------------------------------
# AWS Load Balancer Controller
# watches Gateway/HTTPRoute objects → creates ALB
# -----------------------------------------------
resource "helm_release" "aws_load_balancer_controller" {
  count = var.enable_gateway_api ? 1 : 0

  name       = "aws-load-balancer-controller"
  repository = "https://aws.github.io/eks-charts"
  chart      = "aws-load-balancer-controller"
  version    = "1.8.1"
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

  depends_on = [helm_release.gateway_api[0]]
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

  depends_on = [helm_release.aws_load_balancer_controller[0]]
}

resource "null_resource" "lbc_gateway_crds" {
  count = var.enable_gateway_api ? 1 : 0
  provisioner "local-exec" {
    command = "kubectl apply -f https://raw.githubusercontent.com/kubernetes-sigs/aws-load-balancer-controller/refs/heads/main/config/crd/gateway/gateway-crds.yaml"
  }

  depends_on = [helm_release.aws_load_balancer_controller[0]]
}

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

  depends_on = [
    helm_release.aws_load_balancer_controller[0],
    helm_release.gateway_api[0],
  ]
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


resource "kubectl_manifest" "gateway_api_crds" {
  count = var.enable_gateway_api ? 1 : 0
  
  yaml_body = file("https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.0.0/standard-install.yaml")
  
  depends_on = [aws_eks_cluster.medee-cluster]
}

# Alternative: Use Helm chart for Gateway API CRDs (recommended)
resource "helm_release" "gateway_api" {
  count = var.enable_gateway_api ? 1 : 0
  
  name       = "gateway-api"
  repository = "https://charts.gateway-api.sigs.k8s.io"
  chart      = "gateway-api"
  version    = "1.0.0"
  namespace  = "gateway-system"
  
  create_namespace = true
  
  values = [
    yamlencode({
      image = {
        repository = "registry.k8s.io/gateway-api/admission-server"
        tag        = "v1.0.0"
      }
    })
  ]
  
  depends_on = [aws_eks_node_group.medee-node-group]
}

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

  values = [
    yamlencode({
      clusterName = var.cluster_name
      region      = var.aws_region

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

      # Replica count for HA
      replicaCount = 2

      # Pod disruption budget
      podDisruptionBudget = {
        maxUnavailable = 1
      }
    })
  ]

  depends_on = [
    aws_iam_role_policy_attachment.aws_load_balancer_controller[0],
    helm_release.gateway_api[0]
  ]
}

resource "helm_release" "cert_manager" {
  count = var.enable_cert_manager ? 1 : 0

  name       = "cert-manager"
  repository = "https://charts.jetstack.io"
  chart      = "cert-manager"
  version    = "v1.15.3"
  namespace  = "cert-manager"

  create_namespace = true
  
  wait            = true
  wait_for_jobs   = true
  timeout         = 600
  atomic          = true
  cleanup_on_fail = true

  values = [
    yamlencode({
      installCRDs = true

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

      # Enable prometheus monitoring
      prometheus = {
        enabled = true
        servicemonitor = {
          enabled = true
        }
      }
    })
  ]

  depends_on = [aws_eks_node_group.medee-node-group]
}

resource "kubernetes_manifest" "letsencrypt_cluster_issuer" {
  count = var.enable_cert_manager ? 1 : 0

  manifest = {
    apiVersion = "cert-manager.io/v1"
    kind       = "ClusterIssuer"
    metadata = {
      name = "letsencrypt-prod"
      labels = {
        "app.kubernetes.io/managed-by" = "terraform"
      }
    }
    spec = {
      acme = {
        server = "https://acme-v02.api.letsencrypt.org/directory"
        email  = var.cert_manager_email
        privateKeySecretRef = {
          name = "letsencrypt-prod"
        }
        solvers = [
          {
            http01 = {
              gatewayHTTPRoute = {
                parentRefs = [
                  {
                    name      = "main-gateway"
                    namespace = "gateway-system"
                  }
                ]
              }
            }
          }
        ]
      }
    }
  }

  depends_on = [
    helm_release.cert_manager[0],
    kubernetes_manifest.main_gateway[0]
  ]
}

resource "kubernetes_manifest" "letsencrypt_staging_cluster_issuer" {
  count = var.enable_cert_manager ? 1 : 0

  manifest = {
    apiVersion = "cert-manager.io/v1"
    kind       = "ClusterIssuer"
    metadata = {
      name = "letsencrypt-staging"
      labels = {
        "app.kubernetes.io/managed-by" = "terraform"
      }
    }
    spec = {
      acme = {
        server = "https://acme-staging-v02.api.letsencrypt.org/directory"
        email  = var.cert_manager_email
        privateKeySecretRef = {
          name = "letsencrypt-staging"
        }
        solvers = [
          {
            http01 = {
              gatewayHTTPRoute = {
                parentRefs = [
                  {
                    name      = "main-gateway"
                    namespace = "gateway-system"
                  }
                ]
              }
            }
          }
        ]
      }
    }
  }

  depends_on = [
    helm_release.cert_manager[0],
    kubernetes_manifest.main_gateway[0]
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

resource "kubernetes_manifest" "gateway_class" {
  count = var.enable_gateway_api ? 1 : 0

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
    helm_release.gateway_api[0]
  ]
}

resource "kubernetes_manifest" "main_gateway" {
  count = var.enable_gateway_api ? 1 : 0

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
        # ✅ Added: SSL redirect
        "alb.ingress.kubernetes.io/ssl-redirect" = "443"
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
    kubernetes_namespace.gateway_system[0]
  ]
}

resource "kubernetes_manifest" "main_gateway_certificate" {
  count = var.enable_cert_manager && var.enable_gateway_api ? 1 : 0

  manifest = {
    apiVersion = "cert-manager.io/v1"
    kind       = "Certificate"
    metadata = {
      name      = "main-gateway-tls"
      namespace = "gateway-system"
      labels = {
        "app.kubernetes.io/managed-by" = "terraform"
      }
    }
    spec = {
      secretName = "main-gateway-tls"
      issuerRef = {
        name = var.environment == "prod" ? "letsencrypt-prod" : "letsencrypt-staging"
        kind = "ClusterIssuer"
      }
      dnsNames = [
        "*.${var.domain_name}",
        var.domain_name
      ]
    }
  }

  depends_on = [
    kubernetes_manifest.letsencrypt_cluster_issuer[0],
    kubernetes_namespace.gateway_system[0]
  ]
}

resource "kubernetes_manifest" "reference_grants" {
  for_each = var.enable_gateway_api ? toset([
    "argocd",
    var.monitoring_namespace
  ]) : toset([])

  manifest = {
    apiVersion = "gateway.networking.k8s.io/v1beta1"
    kind       = "ReferenceGrant"
    metadata = {
      name      = "gateway-reference-grant"
      namespace = each.value
      labels = {
        "app.kubernetes.io/managed-by" = "terraform"
      }
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
    kubernetes_namespace.infrastructure_namespaces
  ]
}
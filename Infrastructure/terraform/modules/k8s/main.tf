locals {
  common_tags = {
    Environment = var.environment
    Project     = var.project_name
    ManagedBy   = "terraform"
    Component   = "argocd"
  }
}

# -----------------------------------------------
# ReferenceGrants
# allows HTTPRoutes in these namespaces to attach
# to main-gateway which lives in gateway-system
#
# medee-dev     → frontend + backend app routes
# medee-staging → staging app routes
# argocd        → argocd UI route
# monitoring    → grafana route
# -----------------------------------------------
resource "kubernetes_manifest" "reference_grants" {
  for_each = var.enable_gateway_api ? toset([
    "medee-dev",
    "medee-staging",
    "argocd",
    var.monitoring_namespace,
  ]) : toset([])

  manifest = {
    apiVersion = "gateway.networking.k8s.io/v1beta1"
    kind       = "ReferenceGrant"
    metadata = {
      name      = "gateway-reference-grant"
      namespace = each.value
      labels    = { "app.kubernetes.io/managed-by" = "terraform" }
    }
    spec = {
      from = [{ group = "gateway.networking.k8s.io", kind = "HTTPRoute", namespace = each.value }]
      to   = [{ group = "gateway.networking.k8s.io", kind = "Gateway" }]
    }
  }

  depends_on = [kubernetes_namespace_v1.infrastructure_namespaces]
}

resource "kubernetes_namespace_v1" "infrastructure_namespaces" {
  for_each = toset(var.infrastructure_namespaces)
  metadata {
    name = each.value

    labels = {
      name = each.value
      "app.kubernetes.io/managed-by" = "terraform"
      environment                   = var.environment
      project                       = var.project_name
      purpose                       = "infrastructure"
    }
  }
}

resource "helm_release" "argocd" {
  name       = "argocd"
  repository = "https://argoproj.github.io/argo-helm"
  chart      = "argo-cd"
  version    = "9.4.16"
  namespace  = "argocd"

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

  # depends_on = [
  #   kubernetes_namespace_v1.infrastructure_namespaces,
  # ]
}

resource "kubernetes_manifest" "root_application" {
  manifest = {
    apiVersion = "argoproj.io/v1alpha1"
    kind       = "Application"
    metadata = {
      name      = "root-app"
      namespace = "argocd"
      finalizers = [
        "resources-finalizer.argocd.argoproj.io"
      ]
    }
    spec = {
      project = "default"
      source = {
        repoURL        = var.gitops_repo_url
        targetRevision = var.gitops_target_revision
        path           = "bootstrap"
      }
      destination = {
        server    = "https://kubernetes.default.svc"
        namespace = "argocd"
      }
      syncPolicy = {
        automated = {
          prune    = true
          selfHeal = true
        }
        syncOptions = [
          "CreateNamespace=true",
          "ServerSidceApply=true"
        ]
      }
    }
  }

  depends_on = [ helm_release.argocd]
}

resource "kubernetes_cluster_role_v1" "argocd_application_controller" {
  metadata {
    name = "argocd-application-controller-cluster-admin"
    labels = merge(local.common_tags, {
      "app.kubernetes.io/name"      = "argocd"
      "app.kubernetes.io/component" = "application-controller"
    })
  }

  rule {
    api_groups = ["*"]
    resources  = ["*"]
    verbs      = ["*"]
  }
}

resource "kubernetes_cluster_role_binding_v1" "argocd_application_controller" {
  metadata {
    name = "argocd-application-controller-cluster-admin"
    labels = merge(local.common_tags, {
      "app.kubernetes.io/name"      = "argocd"
      "app.kubernetes.io/component" = "application-controller"
    })
  }

  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = kubernetes_cluster_role_v1.argocd_application_controller.metadata[0].name
  }

  subject {
    kind      = "ServiceAccount"
    name      = "argocd-application-controller"
    namespace = "argocd"
  }

  depends_on = [helm_release.argocd]
}

# in k8s.tf or separate vault.tf
resource "helm_release" "vault" {
  name       = "vault"
  repository = "https://helm.releases.hashicorp.com"
  chart      = "vault"
  namespace  = "vault"

  values = [yamlencode({
    server = {
      ha = { enabled = false }   # single node for dev
    }
    injector = {
      enabled = true             # this is the sidecar injector
    }
  })]

  create_namespace = true
  # depends_on = [ aws_eks_node_group.medee-node-group ]
}
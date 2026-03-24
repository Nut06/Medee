data "aws_eks_cluster" "medee-cluster" {
  name = var.cluster_name
}

data "aws_eks_auth" "medee-cluster" {
  name = var.cluster_name
}

locals {
  common_tags = {
    Environment = var.environment
    Project     = var.project_name
    ManagedBy   = "terraform"
    Component   = "argocd"
  }
}

provider "kubernetes" {
  config_path            = "~/.kube/config"
  config_context         = "my-context"
  host                   = data.aws_eks_cluster.medee-cluster.endpoint
  cluster_ca_certificate = base64decode(data.aws_eks_cluster.medee-cluster.certificate_authority[0].data)
  exec {
    api_version = "client.authentication.k8s.io/v1"
    command     = "aws"
    args        = [
      "eks", "get-token",
      "--cluster-name", var.cluster_name,
      "--region", var.aws_region
    ]
  }
}

provider "helm" {
  kubernetes = {
    host                   = data.aws_eks_cluster.medee-cluster.endpoint
    cluster_ca_certificate = base64decode(data.aws_eks_cluster.medee-cluster.certificate_authority[0].data)
    exec = {
      api_version = "client.authentication.k8s.io/v1"
      command     = "aws"
      args        = [
                      "eks", "get-token",
                      "--cluster-name", var.cluster_name,
                      "--region", var.aws_region
                    ]
    }
  }
}

resource "kubernetes_namespace" "infrastructure_namespaces" {
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

  depends_on = [ aws_eks_cluster.medee-cluster ]
}

resource "helm_release" "argocd" {
  name       = "argocd"
  repository = "https://argoproj.github.io/argo-helm"
  chart      = "argo-cd"
  version    = "7.8.10"
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

  depends_on = [
    kubernetes_namespace.infrastructure_namespaces,
    aws_eks_node_group.medee-node-group
  ]
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

  depends_on = [helm_release.argocd]
}
resource "kubernetes_cluster_role" "argocd_application_controller" {
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

resource "kubernetes_cluster_role_binding" "argocd_application_controller" {
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
    name      = kubernetes_cluster_role.argocd_application_controller.metadata[0].name
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
}

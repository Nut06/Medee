locals {
  common_tags = {
    Environment = var.environment
    Project     = var.project_name
    ManagedBy   = "terraform"
    Component   = "argocd"
  }
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

  # depends_on = [helm_release.argocd]
}

resource "kubernetes_namespace_v1" "medee_dev" {

  metadata {
    name = "medee-dev"
    labels = {
      name                           = "medee-dev"
      # "app.kubernetes.io/managed-by" = "terraform"
      environment                    = var.environment
      project                        = var.project_name
    }
  }

  # depends_on = [helm_release.aws_load_balancer_controller[0]]
}
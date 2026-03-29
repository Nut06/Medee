output "argocd_namespace" {
  value = "argocd"
}

output "infrastructure_namespaces" {
  value = [for ns in kubernetes_namespace_v1.infrastructure_namespaces : ns.metadata[0].name]
}

# output "argocd_root_app_name" {
#   value = kubernetes_manifest.root_application.manifest.metadata.name
# }

# output "vault_release_name" {
#   value = helm_release.vault.name
# }

output "argocd_initial_password_command" {
  value = "kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath=\"{.data.password}\" | base64 -d; echo"
  description = "Command to get ArgoCD initial admin password"
}
# modules/gateway/output.tf

output "gateway_class_name" {
  value = var.enable_gateway_api ? var.gateway_class_name : null
}

output "main_gateway_name" {
  value = var.enable_gateway_api ? "main-gateway" : null
}

output "gateway_namespace" {
  value = var.enable_gateway_api ? "gateway-system" : null
}

output "application_urls" {
  value = var.enable_gateway_api ? {
    argocd   = "https://argocd.${var.domain_name}"
    grafana  = "https://grafana.${var.domain_name}"
    frontend = "https://frontend.${var.domain_name}"
    backend  = "https://backend.${var.domain_name}"
  } : {}
}

output "get_alb_hostname_command" {
  value = var.enable_gateway_api ? "kubectl get gateway main-gateway -n gateway-system -o jsonpath='{.status.addresses[0].value}'" : null
}

# output "cert_manager_release_name" {
#   value = var.enable_cert_manager ? helm_release.cert_manager[0].name : null
# }

output "alb_controller_release_name" {
  value = var.enable_gateway_api ? helm_release.aws_load_balancer_controller[0].name : null
}

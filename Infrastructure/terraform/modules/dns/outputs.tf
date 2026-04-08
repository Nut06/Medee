output "alb_hostname" {
  description = "ALB hostname from Gateway service"
  value       = data.kubernetes_service.alb_hostname.status[0].load_balancer[0].ingress[0].hostname
}

output "root_domain" {
  description = "Root domain record"
  value       = aws_route53_record.root.fqdn
}

output "app_domain" {
  description = "App subdomain record"
  value       = aws_route53_record.app.fqdn
}

output "api_domain" {
  description = "API subdomain record"
  value       = aws_route53_record.api.fqdn
}

output "argocd_domain" {
  description = "ArgoCD subdomain record"
  value       = aws_route53_record.argocd.fqdn
}

output "grafana_domain" {
  description = "Grafana subdomain record"
  value       = aws_route53_record.grafana.fqdn
}

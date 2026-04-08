# -----------------------------------------------
# DNS Module
# จัดการ Route53 DNS Records ทั้งหมด
# -----------------------------------------------

# Data source: ดึง ALB hostname จาก Gateway
data "kubernetes_service" "alb_hostname" {
  metadata {
    name      = var.gateway_service_name
    namespace = var.gateway_namespace
  }
}

# Root domain (dodee.me)
resource "aws_route53_record" "root" {
  zone_id = var.route53_zone_id
  name    = var.domain_name
  type    = "CNAME"
  ttl     = var.dns_ttl
  records = [data.kubernetes_service.alb_hostname.status[0].load_balancer[0].ingress[0].hostname]
}

# App subdomain (app.dodee.me)
resource "aws_route53_record" "app" {
  zone_id = var.route53_zone_id
  name    = "app.${var.domain_name}"
  type    = "CNAME"
  ttl     = var.dns_ttl
  records = [data.kubernetes_service.alb_hostname.status[0].load_balancer[0].ingress[0].hostname]
}

# API subdomain (api.dodee.me)
resource "aws_route53_record" "api" {
  zone_id = var.route53_zone_id
  name    = "api.${var.domain_name}"
  type    = "CNAME"
  ttl     = var.dns_ttl
  records = [data.kubernetes_service.alb_hostname.status[0].load_balancer[0].ingress[0].hostname]
}

# ArgoCD subdomain (argocd.dodee.me)
resource "aws_route53_record" "argocd" {
  zone_id = var.route53_zone_id
  name    = "argocd.${var.domain_name}"
  type    = "CNAME"
  ttl     = var.dns_ttl
  records = [data.kubernetes_service.alb_hostname.status[0].load_balancer[0].ingress[0].hostname]
}

# Grafana subdomain (grafana.dodee.me)
resource "aws_route53_record" "grafana" {
  zone_id = var.route53_zone_id
  name    = "grafana.${var.domain_name}"
  type    = "CNAME"
  ttl     = var.dns_ttl
  records = [data.kubernetes_service.alb_hostname.status[0].load_balancer[0].ingress[0].hostname]
}

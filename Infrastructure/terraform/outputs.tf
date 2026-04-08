# --- EC2 Instance Outputs ---

output "jenkins_public_ip" {
  description = "The public IP address assigned to the Jenkins instance"
  value       = aws_instance.ec2_jenkins.public_ip
}

output "jenkins_instance_id" {
  description = "The ID of the EC2 instance"
  value       = aws_instance.ec2_jenkins.id
}

output "jenkins_url" {
  description = "The URL to access the Jenkins Web UI"
  value       = "http://${aws_instance.ec2_jenkins.public_ip}:8080"
}

# --- VPC & Networking Outputs ---

output "vpc_id" {
  description = "The ID of the VPC"
  value       = module.vpc.vpc_id
}

output "public_subnets" {
  description = "List of IDs of public subnets"
  value       = module.vpc.public_subnets
}

# --- Security Group Outputs ---

output "jenkins_sg_id" {
  description = "The ID of the security group created for Jenkins"
  value       = module.jenkins_sg.security_group_id
}

# ... existing outputs ...

output "gateway_class_name" {
  description = "Gateway class name"
  value       = var.enable_gateway_api ? var.gateway_class_name : null
}

output "main_gateway_name" {
  description = "Main gateway name"
  value       = var.enable_gateway_api ? "main-gateway" : null
}

# # ACM Certificate
output "acm_certificate_arn" {
  description = "ARN of the validated ACM certificate"
  value       = aws_acm_certificate.main.arn
}

# output "application_urls" {
#   description = "Application URLs via Gateway API"
#   value = var.enable_gateway_api ? {
#     grafana = "https://grafana.${var.domain_name}"
#     argocd  = "https://argocd.${var.domain_name}"
#     frontend = "https://frontend.${var.domain_name}"
#     backend = "https://backend.${var.domain_name}"
#   } : {}
# }

# output "gateway_load_balancer_command" {
#   description = "Command to get Gateway LoadBalancer hostname"
#   value = var.enable_gateway_api ? "kubectl get gateway main-gateway -n gateway-system -o jsonpath='{.status.addresses[0].value}'" : null
# }

# output "argocd_initial_password_command" {
#   value = module.k8s.argocd_initial_password_command
# }

# output "route53_nameservers" {
#   description = "Route53 Name Servers (update at domain registrar)"
#   value = aws_route53_zone.main.name_servers
# }

# # Route53
# output "route53_zone_id" {
#   description = "Route53 Hosted Zone ID"
#   value       = aws_route53_zone.main.zone_id
# }

# # ALB
# output "alb_hostname" {
#   description = "ALB hostname"
#   value       = module.dns.alb_hostname
# }

# # DNS Records
# output "root_domain" {
#   description = "Root domain"
#   value       = module.dns.root_domain
# }

# output "app_domain" {
#   description = "App subdomain"
#   value       = module.dns.app_domain
# }

# output "api_domain" {
#   description = "API subdomain"
#   value       = module.dns.api_domain
# }

# output "argocd_domain" {
#   description = "ArgoCD subdomain"
#   value       = module.dns.argocd_domain
# }

# output "grafana_domain" {
#   description = "Grafana subdomain"
#   value       = module.dns.grafana_domain
# }

# # Application URLs
# output "root_url" {
#   description = "Root domain URL"
#   value       = "https://${var.domain_name}"
# }

# output "app_url" {
#   description = "Application URL"
#   value       = "https://app.${var.domain_name}"
# }

# output "api_url" {
#   description = "API URL"
#   value       = "https://api.${var.domain_name}"
# }

# output "argocd_url" {
#   description = "ArgoCD URL"
#   value       = "https://argocd.${var.domain_name}"
# }

# output "grafana_url" {
#   description = "Grafana URL"
#   value       = "https://grafana.${var.domain_name}"
# }
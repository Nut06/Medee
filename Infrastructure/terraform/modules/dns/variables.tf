variable "route53_zone_id" {
  description = "Route53 Hosted Zone ID"
  type        = string
}

variable "domain_name" {
  description = "Domain name"
  type        = string
}

variable "gateway_service_name" {
  description = "Gateway service name in Kubernetes"
  type        = string
  default     = "main-gateway"
}

variable "gateway_namespace" {
  description = "Gateway namespace in Kubernetes"
  type        = string
  default     = "gateway-system"
}

variable "dns_ttl" {
  description = "DNS TTL in seconds"
  type        = number
  default     = 300
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "project_name" {
  description = "Project name"
  type        = string
}

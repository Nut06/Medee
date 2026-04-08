variable "enable_gateway_api" {
  description = "Enable Gateway API instead of Ingress"
  type        = bool
  default     = true
}

variable "cluster_name" {
  type = string
}

variable "aws_region" {
  type = string
}
variable "enable_cert_manager" {
  description = "Enable cert-manager for TLS certificates"
  type        = bool
  default     = true
}

variable "cert_manager_email" {
  description = "Email for Let's Encrypt certificates"
  type        = string
}

variable "gateway_class_name" {
  description = "Gateway class name to use"
  type        = string
  default     = "aws-load-balancer-controller"
}

variable "domain_name" {
  description = "Domain name for applications"
  type        = string
  default     = "medee.local"
}

variable "environment" {
  type = string
}

variable "project_name" {
  default = "medee"
    type = string
}

variable "lb_controller_role_arn" {
  description = "IAM role ARN for AWS Load Balancer Controller (IRSA)"
  type        = string
}

variable "acm_certificate_arn" {
  description = "ACM certificate ARN for HTTPS on the ALB"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID where the EKS cluster is deployed"
  type        = string
}

variable "public_subnet_ids" {
  description = "Public subnet IDs for ALB placement"
  type        = list(string)
}
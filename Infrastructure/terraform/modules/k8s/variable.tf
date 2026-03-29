variable "cluster_name" {
  type = string
}

variable "environment" {
  type = string
}

variable "project_name" {
  default = "medee"
    type = string
}

variable "aws_region" {
  type = string
}

variable "infrastructure_namespaces" {
  description = "Infrastructure namespaces managed by Terraform"
  type        = list(string)
  default     = [
    "argocd",
    "gateway-system"
  ]
}

variable "gitops_repo_url" {
  description = "GitOps repository URL"
  type        = string
}

variable "gitops_target_revision" {
  description = "GitOps repository target revision"
  type        = string
  default     = "HEAD"
}

variable "enable_gateway_api" {
  description = "Enable Gateway API instead of Ingress"
  type        = bool
  default     = true
}

variable "monitoring_namespace" {
  description = "Monitoring namespace"
  type        = string
  default     = "monitoring"
}
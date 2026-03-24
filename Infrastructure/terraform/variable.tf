variable "my_key_name" {
  type = string
}

variable "azs" {
  default = [ "ap-southeast-7a" ]
  type = list(string)
}

variable "public_subnets" {
    default = [ "10.0.101.0/24" ]
  type = list(string)
}

variable "project_name" {
  default = "medee"
    type = string
}

variable "cluster_name" {
  type = string
}

variable "aws_region" {
  type = string
}

variable "environment" {
  type = string
}

variable "argocd_namespace" {
  description = "Namespace for ArgoCD"
  type        = string
  default     = "argocd"
}

# ... existing variables ...
variable "gateway_class_name" {
  description = "Gateway class name to use"
  type        = string
  default     = "aws-load-balancer-controller"
}

variable "enable_gateway_api" {
  description = "Enable Gateway API instead of Ingress"
  type        = bool
  default     = true
}

variable "domain_name" {
  description = "Domain name for applications"
  type        = string
  default     = "medee.local"
}

variable "enable_cert_manager" {
  description = "Enable cert-manager for TLS certificates"
  type        = bool
  default     = true
}

variable "enable_external_dns" {
  description = "Enable external-dns for automatic DNS management"
  type        = bool
  default     = false
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

variable "fargate_namespaces" {
  description = "Namespaces for Fargate profiles"
  type        = list(string)
  default     = ["frontend-dev", "backend-dev"]
}

variable "node_group_instance_types" {
  description = "Instance types for the node group"
  type        = list(string)
  default     = ["t3.medium", "t3.large"]
}

variable "node_group_desired_size" {
  description = "Desired number of nodes"
  type        = number
  default     = 2
}

variable "node_group_max_size" {
  description = "Maximum number of nodes"
  type        = number
  default     = 10
}

variable "node_group_min_size" {
  description = "Minimum number of nodes"
  type        = number
  default     = 1
}



variable "fargate_namespaces" {
  description = "Namespaces that should run on Fargate"
  type        = list(string)
  default     = [
    "monitoring",      # Run monitoring on Fargate for cost optimization
    "frontend-dev",
    "backend-dev"
  ]
}

variable "infrastructure_namespaces" {
  description = "Infrastructure namespaces managed by Terraform"
  type        = list(string)
  default     = [
    "argocd",
    "gateway-system"
  ]
}

variable "cert_manager_email" {
  description = "Email for Let's Encrypt certificates"
  type        = string
}

variable "monitoring_namespace" {
  description = "Monitoring namespace"
  type        = string
  default     = "monitoring"
}

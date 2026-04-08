
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

variable "node_group_instance_types" {
  description = "Instance types for the node group"
  type        = list(string)
  default     = ["t3.medium", "t3.large"]
}

variable "node_group_desired_size" {
  description = "Desired size for the node group"
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

variable "enable_gateway_api" {
  description = "Enable Gateway API instead of Ingress"
  type        = bool
  default     = true
}

variable "private_subnets" {
  type = list(string)
}

variable "public_subnets" {
  type = list(string)
}

variable "domain_name" {
  type = string
}
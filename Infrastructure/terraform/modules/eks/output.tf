output "cluster_name" {
  value = aws_eks_cluster.medee-cluster.name
}

output "cluster_endpoint" {
  value = aws_eks_cluster.medee-cluster.endpoint
}

output "cluster_ca" {
  value = aws_eks_cluster.medee-cluster.certificate_authority[0].data
}

output "cluster_oidc_issuer" {
  value = aws_eks_cluster.medee-cluster.identity[0].oidc[0].issuer
}

output "oidc_provider_arn" {
  value = aws_iam_openid_connect_provider.eks[0].arn
}

output "lb_controller_role_arn" {
  value = var.enable_gateway_api ? aws_iam_role.aws_load_balancer_controller[0].arn : null
}

output "node_group_id" {
  value = aws_eks_node_group.medee-node-group.id
}

output "fargate_profile_ids" {
  value = { for k, v in aws_eks_fargate_profile.medee-fargate : k => v.id }
}
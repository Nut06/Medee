# EKS Cluster
resource "aws_eks_cluster" "medee-cluster" {
  name     = var.cluster_name
  role_arn = aws_iam_role.eks-cluster-role.arn
  version  = "1.31"

  enabled_cluster_log_types = ["api", "audit", "authenticator", "controllerManager", "scheduler"]

  vpc_config {
    endpoint_private_access = false 
    endpoint_public_access  = true
    
    subnet_ids = concat(
      var.private_subnets,
      var.public_subnets
    )
  }

  access_config {
    authentication_mode = "API"
    bootstrap_cluster_creator_admin_permissions = true
  }

  depends_on = [
    aws_iam_role_policy_attachment.cluster_EKSClusterPolicy,
  ]

  tags = {
    Name        = var.cluster_name
    Environment = var.environment
    Project     = var.project_name
  }
}

resource "aws_eks_node_group" "medee-node-group" {
  cluster_name    = aws_eks_cluster.medee-cluster.name
  node_group_name = "${var.project_name}-node-group"
  node_role_arn   = aws_iam_role.eks-node-role.arn
  subnet_ids      = var.private_subnets

  capacity_type  = "ON_DEMAND"
  instance_types = var.node_group_instance_types

  scaling_config {
    desired_size = var.node_group_desired_size
    max_size     = var.node_group_max_size
    min_size     = var.node_group_min_size
  }

  update_config {
    max_unavailable = 1
  }

  depends_on = [
    aws_iam_role_policy_attachment.node_AmazonEKSWorkerNodeMinimalPolicy,
    aws_iam_role_policy_attachment.node_AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.node_AmazonEC2ContainerRegistryReadOnly,
  ]

  tags = {
    Name        = "${var.project_name}-node-group"
    Environment = var.environment
    Project     = var.project_name
  }
}

resource "aws_eks_fargate_profile" "medee-fargate" {
  for_each = toset(var.fargate_namespaces)
  
  cluster_name           = aws_eks_cluster.medee-cluster.name
  fargate_profile_name   = "${var.project_name}-fargate-${each.value}"
  pod_execution_role_arn = aws_iam_role.eks-fargate-role.arn
  subnet_ids             = toset(var.private_subnets)

  selector {
    namespace = each.value
    labels = {
      "compute-type" = "fargate"
    }
  }

  depends_on = [
    aws_iam_role_policy_attachment.fargate_AmazonEKSFargatePodExecutionRolePolicy,
  ]

  tags = {
    Name        = "${var.project_name}-fargate-${each.value}"
    Environment = var.environment
    Project     = var.project_name
  }
}

# EKS Add-ons
resource "aws_eks_addon" "vpc_cni" {
  cluster_name = aws_eks_cluster.medee-cluster.name
  addon_name   = "vpc-cni"
  
  depends_on = [aws_eks_node_group.medee-node-group]
}

resource "aws_eks_addon" "coredns" {
  cluster_name = aws_eks_cluster.medee-cluster.name
  addon_name   = "coredns"
  
  depends_on = [aws_eks_node_group.medee-node-group]
}

resource "aws_eks_addon" "kube_proxy" {
  cluster_name = aws_eks_cluster.medee-cluster.name
  addon_name   = "kube-proxy"
  
  depends_on = [aws_eks_node_group.medee-node-group]
}

resource "aws_eks_addon" "ebs_csi" {
  cluster_name = aws_eks_cluster.medee-cluster.name
  addon_name   = "aws-ebs-csi-driver"
  service_account_role_arn = aws_iam_role.ebs_csi.arn
  depends_on = [aws_eks_node_group.medee-node-group]
}
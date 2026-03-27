
module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name                    = "${var.project_name}-vpc"
  cidr                    = "10.0.0.0/16"
  map_public_ip_on_launch = true

  azs             = var.azs
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24"]
  public_subnets  = var.public_subnets

  enable_nat_gateway   = true
  single_nat_gateway = true
  enable_vpn_gateway   = false
  enable_dns_hostnames = true # ← เพิ่มนี่
  enable_dns_support   = true

  #   public_subnet_tags = {
  #     # "kubernetes.io/role/elb" = "1"
  #   }
}
module "jenkins_sg" {
  source      = "terraform-aws-modules/security-group/aws"
  name        = "jenkins-service-sg"
  description = "Security for jenkins web UI"
  vpc_id      = module.vpc.vpc_id

  ingress_cidr_blocks = ["0.0.0.0/0"]

  ingress_rules = ["ssh-tcp"]

  ingress_with_cidr_blocks = [
    {
      description = "Jenkins Web UI"
      from_port   = 8080
      to_port     = 8080
      protocol    = "tcp"
      cidr_blocks = "0.0.0.0/0"
    },
    {
      description = "Enable for Sonarqube"
      from_port   = 9000
      to_port     = 9000
      protocol    = "tcp"
      cidr_blocks = "0.0.0.0/0"
    }
  ]

  egress_rules = ["all-all"]
}

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd-gp3/ubuntu-noble-24.04-amd64-server-*"]
  }

  owners = ["099720109477"] # Canonical
}

# terraform destroy -target=module.ec2_jenkins
# ec2_jenkins
# module "ec2_jenkins" {
#   source  = "terraform-aws-modules/ec2-instance/aws"

#   associate_public_ip_address = true
#   name = "${var.project_name}-jenkins"
#   ami = data.aws_ami.ubuntu.id
#   instance_type = "t3.large"
#   key_name      = aws_key_pair.jenkins_key.id
#   monitoring    = true
#   subnet_id     = module.vpc.public_subnets[0]

#   root_block_device = {
#       volume_type = "gp3"
#       volume_size = 30
#       delete_on_termination = true
#     }

#   vpc_security_group_ids = [module.jenkins_sg.security_group_id]
#   tags = {
#     Project = var.project_name
#   }
# }

data "aws_eks_cluster" "medee-cluster" {
  name       = var.cluster_name
  depends_on = [module.eks]
}

resource "aws_key_pair" "jenkins_key" {
  key_name   = var.my_key_name
  public_key = file("~/.ssh/jenkins-ec2-key.pub")
}

module "eks" {
  source                    = "./modules/eks"
  cluster_name              = var.cluster_name
  environment               = var.environment
  project_name              = var.project_name
  node_group_instance_types = var.node_group_instance_types
  node_group_desired_size   = var.node_group_desired_size
  node_group_max_size       = var.node_group_max_size
  node_group_min_size       = var.node_group_min_size
  fargate_namespaces        = var.fargate_namespaces
  enable_gateway_api        = var.enable_gateway_api
  private_subnets = module.vpc.private_subnets
  public_subnets = module.vpc.public_subnets

  depends_on                = [module.vpc, module.jenkins_sg]
}
module "gateway" {
  source              = "./modules/gateway"
  enable_gateway_api  = var.enable_gateway_api
  cluster_name        = var.cluster_name
  aws_region          = var.aws_region
  enable_cert_manager = var.enable_cert_manager

  cert_manager_email = var.cert_manager_email
  public_subnet_ids = module.vpc.public_subnets
  gateway_class_name = var.gateway_class_name
  domain_name        = var.domain_name
  lb_controller_role_arn = module.eks.lb_controller_role_arn

  monitoring_namespace = var.monitoring_namespace
  vpc_id = module.vpc.vpc_id
  environment = var.environment
  project_name = var.project_name
  acm_certificate_arn = var.acm_certificate_arn

  # depends_on   = [module.k8s]
}


module "k8s" {
  source       = "./modules/k8s"
  cluster_name = var.cluster_name
  environment  = var.environment
  project_name = var.project_name

  enable_gateway_api = var.enable_gateway_api
  monitoring_namespace = var.monitoring_namespace
  aws_region                = var.aws_region
  infrastructure_namespaces = var.infrastructure_namespaces
  gitops_repo_url           = var.gitops_repo_url
  gitops_target_revision    = var.gitops_target_revision
  depends_on                = [module.gateway]
}


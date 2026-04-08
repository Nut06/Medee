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


resource "aws_key_pair" "jenkins_key" {
  key_name   = var.my_key_name
  public_key = file("~/.ssh/jenkins-ec2-key.pub")
}

module "eks" {
  source                    = "./modules/eks"
  domain_name = var.domain_name
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

# create acm certificate & validation before create gateway
resource "aws_acm_certificate" "main" {
  domain_name               = var.domain_name
  subject_alternative_names = ["*.${var.domain_name}"]
  validation_method         = "DNS"

  tags = {
    Environment = var.environment
    Project     = var.project_name
    ManagedBy   = "terraform"
}

  lifecycle {
    create_before_destroy = true
  }
}


data "aws_eks_cluster" "medee-cluster" {
  name       = var.cluster_name
  depends_on = [module.eks]
}

resource "aws_route53_zone" "main" {
  name = var.domain_name   # "dodee.me"

  tags = {
    Environment = var.environment
    Project     = var.project_name
  }
}

# ลบ locals block ทิ้ง (บรรทัด 147-157)

# แก้ aws_route53_record ให้ใช้ for_each แทน
# ใช้ตอนที่จะ validation
# resource "aws_route53_record" "cert_validation" {
#   for_each = {
#     for dvo in aws_acm_certificate.main.domain_validation_options : dvo.domain_name => {
#       name   = dvo.resource_record_name
#       record = dvo.resource_record_value
#       type   = dvo.resource_record_type
#     }
#   }

#   allow_overwrite = true
#   name            = each.value.name
#   records         = [each.value.record]
#   ttl             = 60
#   type            = each.value.type
#   zone_id         = aws_route53_zone.main.zone_id
# }

# แก้ validation ให้รอทุก records
# resource "aws_acm_certificate_validation" "main" {
#   certificate_arn         = aws_acm_certificate.main.arn
#   validation_record_fqdns = [for record in aws_route53_record.cert_validation : record.fqdn]
# }


# gateway

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

  vpc_id = module.vpc.vpc_id
  environment = var.environment
  project_name = var.project_name
  acm_certificate_arn = aws_acm_certificate.main.arn

  depends_on = [
    aws_acm_certificate.main
  ]
}

module "k8s" {
  source       = "./modules/k8s"
  cluster_name = var.cluster_name
  environment  = var.environment
  project_name = var.project_name

  enable_gateway_api = var.enable_gateway_api
  aws_region                = var.aws_region
  gitops_repo_url           = var.gitops_repo_url
  gitops_target_revision    = var.gitops_target_revision
  depends_on                = [module.gateway]
}

# module "dns" {
#   source = "./modules/dns"

#   route53_zone_id      = aws_route53_zone.main.zone_id
#   domain_name          = var.domain_name
#   gateway_service_name = "main-gateway"
#   gateway_namespace    = "gateway-system"
#   dns_ttl              = 300
#   environment          = var.environment
#   project_name         = var.project_name

#   depends_on = [
#     module.gateway,
#     module.k8s
#   ]
# }
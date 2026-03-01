provider "aws" {
  region = "ap-southeast-7"
}

  module "vpc" {
    source = "terraform-aws-modules/vpc/aws"

    name = "${var.project_name}-vpc"
    cidr = "10.0.0.0/16"
    map_public_ip_on_launch = true

    azs             = var.azs
  #   private_subnets = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
    public_subnets  = var.public_subnets

    enable_nat_gateway = false
    enable_vpn_gateway = false

  #   public_subnet_tags = {
  #     # "kubernetes.io/role/elb" = "1"
  #   }
  }

  module "jenkins_sg" {
    source = "terraform-aws-modules/security-group/aws"
    name = "jenkins-service-sg"
    description = "Security for jenkins web UI"
    vpc_id = module.vpc.vpc_id

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
    name = "name"
    values = ["ubuntu/images/hvm-ssd-gp3/ubuntu-noble-24.04-amd64-server-*"]
  }

  owners = ["099720109477"] # Canonical
}

  # ec2_jenkins
  module "ec2_jenkins" {
    source  = "terraform-aws-modules/ec2-instance/aws"

    associate_public_ip_address = true
    name = "${var.project_name}-jenkins"
    ami = data.aws_ami.ubuntu.id
    instance_type = "t3.medium"
    key_name      = aws_key_pair.jenkins_key.id
    monitoring    = true
    subnet_id     = module.vpc.public_subnets[0]
    
    root_block_device = {
        volume_type = "gp3"
        volume_size = 14
        delete_on_termination = true
      }

    user_data = <<-EOF
                  #!/bin/bash
                  sudo yum update -y
                  sudo dnf install -y docker
                  sudo systemctl enable --now docker
                  EOF

    vpc_security_group_ids = [module.jenkins_sg.security_group_id]
    tags = {
      Project = var.project_name
    }
  }

  resource "aws_key_pair" "jenkins_key" {
    key_name = var.my_key_name
    public_key = file("~/.ssh/jenkins-ec2-key.pub")
  }
# --- EC2 Instance Outputs ---

output "jenkins_public_ip" {
  description = "The public IP address assigned to the Jenkins instance"
  value       = module.ec2_jenkins.public_ip
}

output "jenkins_instance_id" {
  description = "The ID of the EC2 instance"
  value       = module.ec2_jenkins.id
}

output "jenkins_url" {
  description = "The URL to access the Jenkins Web UI"
  value       = "http://${module.ec2_jenkins.public_ip}:8080"
}

# --- VPC & Networking Outputs ---

output "vpc_id" {
  description = "The ID of the VPC"
  value       = module.vpc.vpc_id
}

output "public_subnets" {
  description = "List of IDs of public subnets"
  value       = module.vpc.public_subnets
}

# --- Security Group Outputs ---

output "jenkins_sg_id" {
  description = "The ID of the security group created for Jenkins"
  value       = module.jenkins_sg.security_group_id
}
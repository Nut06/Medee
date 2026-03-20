resource "aws_instance" "ec2_jenkins" {
  ami = data.aws_ami.ubuntu.id
  instance_type = "t3.large"
  associate_public_ip_address = true
    key_name      = aws_key_pair.jenkins_key.id
    monitoring    = true
    subnet_id     = module.vpc.public_subnets[0]

    ebs_block_device {
      device_name = "/dev/sdf" 
        volume_size = 25
        volume_type = "gp3"
    }

    root_block_device {
        volume_size = 14
        volume_type = "gp3" 
      delete_on_termination = true
    }

    vpc_security_group_ids = [module.jenkins_sg.security_group_id]
    tags = {
        Project = var.project_name
    }
}   
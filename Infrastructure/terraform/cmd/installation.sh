#!/bin/bash

# install
sudo apt update -y
sleep 5
sudo apt install fontconfig openjdk-21-jre -y
java -version

sleep 5
sudo wget -O /etc/apt/keyrings/jenkins-keyring.asc \
  https://pkg.jenkins.io/debian-stable/jenkins.io-2026.key
echo "deb [signed-by=/etc/apt/keyrings/jenkins-keyring.asc]" \
  https://pkg.jenkins.io/debian-stable binary/ | sudo tee \
  /etc/apt/sources.list.d/jenkins.list > /dev/null

sleep 5
sudo apt update -y
sleep 5
sudo apt install jenkins -y

# enable jenkins
sudo systemctl enable jenkins
sudo systemctl start jenkins
sleep 10

# install go
sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt update
sudo apt install golang-go

sleep 10
sudo snap install --classic go

# 2. Install Docker &  SONARQUBE  AS A CONTAINER IN THE EC2 INSTANCE
sudo apt-get update -y
sudo apt-get install docker.io -y
sleep 5
sudo usermod -aG docker ubuntu 
sudo usermod -aG docker jenkins 
newgrp docker
sleep 5
sudo chmod 777 /var/run/docker.sock

# 1. อัปเดตแพ็กเกจและติดตั้งเครื่องมือพื้นฐาน
sudo apt-get update
sudo apt-get install -y ca-certificates curl gnupg lsb-release

# 2. เพิ่ม Docker's official GPG key
sudo mkdir -m 0755 -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg

# 3. ตั้งค่า Repository
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sleep 10

sudo apt-get update
sudo apt-get install -y docker-buildx-plugin

# ตรวจสอบเวอร์ชัน (ถ้าขึ้นเลขเวอร์ชันแปลว่ารอดแล้วครับ)
docker buildx version

# รีสตาร์ท Jenkins เพื่อให้สิทธิ์ใหม่มีผล
sudo systemctl restart jenkins

sleep 10
docker run -d --name sonar -p 9000:9000 \
        -v sonarqube_data:/opt/sonarqube/data \
        -v sonarqube_logs:/opt/sonarqube/logs \
        -v sonarqube_extensions:/opt/sonarqube/extensions \
        dhi.io/sonarqube:25-debian13
sleep 5

# 3. Install trivy ON THE INSTANCE
sudo apt-get install wget apt-transport-https gnupg lsb-release -y
sleep 5
wget -qO - https://aquasecurity.github.io/trivy-repo/deb/public.key | gpg --dearmor | sudo tee /usr/share/keyrings/trivy.gpg > /dev/null
echo "deb [signed-by=/usr/share/keyrings/trivy.gpg] https://aquasecurity.github.io/trivy-repo/deb $(lsb_release -sc) main" | sudo tee -a /etc/apt/sources.list.d/trivy.list
sudo apt-get update -y
sleep 5
sudo apt-get install trivy -y
# install trivy plugin
sudo -u jenkins -H trivy --config /dev/null plugin install scan2html
sleep 5

# 5. Install kubectl
sudo apt update
sleep 5
sudo apt install curl -y
sleep 5
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sleep 5
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
kubectl version --client
which kubectl

sleep 10

# 6. Install AWS CLI 
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
sleep 5
sudo apt-get install unzip -y
sleep 5
unzip awscliv2.zip
sudo ./aws/install

#7. eksctl installation 
curl --silent --location "https://github.com/weaveworks/eksctl/releases/latest/download/eksctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp
sleep 5
sudo mv /tmp/eksctl /usr/local/bin
sleep 5

#Status of installation
eksctl version
kubectl version --client
terraform -v
docker -version
aws --version
java --version
#!/bin/bash
# apply resource
terraform apply -target=module.vpc \
  -target=module.eks \
#   -target=module.jenkins_sg \
#   -target=aws_key_pair.jenkins_key \
#   -target=aws_instance.ec2_jenkins \

# Stage 2 — update kubeconfig then apply everything else
sleep 5
aws eks update-kubeconfig --name medee-cluster --region ap-southeast-7

aws eks describe-cluster \
  --region ap-southeast-7 \
  --name medee-cluster \
  --query 'cluster.arn' \
  --output text

sleep 5
terraform apply

terraform destroy -target=module.eks

terraform state rm module.eks.aws_eks_cluster.medee-cluster \
terraform state rm module.eks.aws_iam_role.eks-cluster-role \ 
terraform state rm module.eks.aws_iam_role_policy_attachment.cluster_EKSClusterPolicy \

terraform plan -target=module.k8s -var-file=../terraform.tfvars
terraform apply -auto-approve plan -var-file=../terraform.tfvars

# check for load balance
aws iam get-role \
  --role-name medee-aws-load-balancer-controller \
  --region ap-southeast-7 \
  --query "Role.Arn"

eksctl get iamserviceaccount --cluster=medee-cluster

{
    "CertificateArn": "arn:aws:acm:ap-southeast-7:652253416828:certificate/b46a3556-5d08-4eea-a994-8a5962dd55d0"
}

aws acm list-certificates \
  --region ap-southeast-7 \
  --query "CertificateSummaryList[*].{ARN:CertificateArn,Domain:DomainName}"

kubectl get pods -n kube-system | grep aws-load-balancer
kubectl describe pod -n kube-system -l app.kubernetes.io/name=aws-load-balancer-controller

kubectl get pods -n kube-system -l app.kubernetes.io/name=aws-load-balancer-controller
kubectl logs -n kube-system -l app.kubernetes.io/name=aws-load-balancer-controller

# resource "kubernetes_manifest" "gateway_class" { ... }
# resource "kubernetes_manifest" "main_gateway_lbconfig" { ... }
# resource "kubernetes_manifest" "main_gateway" { ... }

# import the missing resource state
terraform import module.gateway.kubernetes_namespace_v1.gateway_system[0] gateway-system

# 
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d; echo

# port forwarding
kubectl port-forward svc/argocd-server -n argocd 8080:80

terraform apply -target=module.k8s.kubernetes_manifest.root_application
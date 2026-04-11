# Terraform Destroy Guide

คู่มือการลบ Infrastructure ของ Medee Project อย่างถูกต้องและปลอดภัย

---

## ⚠️ คำเตือน

- **สำรองข้อมูลสำคัญก่อนทำการลบ**
- การลบจะไม่สามารถย้อนกลับได้
- ตรวจสอบว่าคุณอยู่ใน environment ที่ถูกต้อง

---

## วิธีที่ 1: Destroy แบบอัตโนมัติ (แนะนำ)

### ขั้นตอนที่ 1: ลบ ArgoCD Applications

```bash
# เข้าไปที่ directory terraform
cd Infrastructure/terraform

# ลบ applications ทั้งหมดที่ ArgoCD deploy ไว้
kubectl delete application --all -n argocd

# รอให้ resources ถูกลบจริงๆ (ประมาณ 2-3 นาที)
kubectl get application -n argocd
```

### ขั้นตอนที่ 2: Destroy ทั้งหมด

```bash
# Destroy ทั้งหมด
terraform destroy

# หรือถ้าต้องการ auto-approve
terraform destroy -auto-approve
```

---

## วิธีที่ 2: Destroy แบบทีละขั้นตอน (ถ้าวิธีที่ 1 มีปัญหา)

### ขั้นตอนที่ 1: ลบ Route53 Records

```bash
cd Infrastructure/terraform

terraform destroy -target=aws_route53_record.app
terraform destroy -target=aws_route53_record.api
terraform destroy -target=aws_route53_record.argocd
terraform destroy -target=aws_route53_record.grafana
```

### ขั้นตอนที่ 2: ลบ ArgoCD Applications

```bash
# ลบ root-app และ applications อื่นๆ
kubectl delete application root-app -n argocd
kubectl delete application --all -n argocd

# ตรวจสอบว่าลบหมดแล้ว
kubectl get application -n argocd
```

### ขั้นตอนที่ 3: ลบ K8s Module

```bash
# ลบ ArgoCD, Vault, Namespaces
terraform destroy -target=module.k8s
```

### ขั้นตอนที่ 4: ลบ Gateway Module

```bash
# ลบ AWS Load Balancer Controller, Gateway API, main-gateway
terraform destroy -target=module.gateway
```

### ขั้นตอนที่ 5: ลบ EKS Cluster

```bash
# ลบ EKS cluster, node groups, fargate profiles
terraform destroy -target=module.eks
```

### ขั้นตอนที่ 6: ลบ Certificate และ Validation

```bash
terraform destroy -target=aws_acm_certificate_validation.main -auto-approve
terraform destroy -target=aws_route53_record.cert_validation -auto-approve
terraform destroy -target=aws_acm_certificate.main -auto-approve
```

### ขั้นตอนที่ 7: ลบ Route53 Zone

```bash
terraform destroy -target=aws_route53_zone.main -auto-approve
```

### ขั้นตอนที่ 8: ลบ Jenkins Resources

```bash
terraform destroy -target=aws_key_pair.jenkins_key -auto-approve
terraform destroy -target=module.jenkins_sg -auto-approve
```

### ขั้นตอนที่ 9: ลบ VPC

```bash
# ลบ VPC, Subnets, NAT Gateway, Internet Gateway
terraform destroy -target=module.vpc -auto-approve
```

### ขั้นตอนที่ 10: ตรวจสอบและลบที่เหลือ

```bash
# ตรวจสอบว่ามี resources เหลืออยู่หรือไม่
terraform state list

# ลบทั้งหมด
terraform destroy
```

---

## การตรวจสอบหลังการลบ

### ตรวจสอบ Terraform State

```bash
# ดู resources ที่เหลือ
terraform state list

# ถ้าไม่มีอะไรเหลือจะไม่แสดงผลอะไร
```

### ตรวจสอบ AWS Console

1. **EC2 Dashboard**
   - Load Balancers (ต้องถูกลบหมด)
   - Target Groups (ต้องถูกลบหมด)

2. **EKS Dashboard**
   - Clusters (ต้องถูกลบหมด)

3. **VPC Dashboard**
   - VPCs (ต้องถูกลบหมด)
   - NAT Gateways (ต้องถูกลบหมด)
   - Elastic IPs (ต้องถูกลบหมด)

4. **Route53 Dashboard**
   - Hosted Zones (ตรวจสอบว่าถูกลบหรือไม่)

5. **Certificate Manager**
   - Certificates (ต้องถูกลบหมด)

---

## แก้ไขปัญหาที่พบบ่อย

### ปัญหา: ALB ไม่ถูกลบ

```bash
# ลบ ALB ด้วยตนเองผ่าน kubectl
kubectl delete gateway main-gateway -n gateway-system

# รอสักครู่แล้วลองอีกครั้ง
terraform destroy -target=module.gateway
```

### ปัญหา: VPC ไม่ถูกลบเพราะมี ENI ค้างอยู่

```bash
# ตรวจสอบ ENI ที่ค้างอยู่
aws ec2 describe-network-interfaces --filters "Name=vpc-id,Values=<VPC_ID>"

# ลบ ENI ด้วยตนเอง
aws ec2 delete-network-interface --network-interface-id <ENI_ID>

# ลอง destroy VPC อีกครั้ง
terraform destroy -target=module.vpc
```

### ปัญหา: Security Group ไม่ถูกลบ

```bash
# ลบ dependencies ก่อน
terraform destroy -target=module.eks
terraform destroy -target=module.gateway

# แล้วค่อยลบ security group
terraform destroy -target=module.jenkins_sg
```

### ปัญหา: PVC ค้างอยู่ใน EKS

```bash
# ลบ PVC ทั้งหมด
kubectl delete pvc --all -n medee-dev
kubectl delete pvc --all -n monitoring
kubectl delete pvc --all -n vault

# รอให้ EBS volumes ถูกลบ
aws ec2 describe-volumes --filters "Name=tag:kubernetes.io/cluster/<CLUSTER_NAME>,Values=owned"
```

---

## ลำดับการ Destroy (สรุป)

```
1. Route53 Records (app, api, argocd, grafana)
   ↓
2. ArgoCD Applications (kubectl delete)
   ↓
3. K8s Module (ArgoCD, Vault, Namespaces)
   ↓
4. Gateway Module (ALB, Gateway API)
   ↓
5. EKS Module (Cluster, Node Groups)
   ↓
6. Certificate Validation & ACM Certificate
   ↓
7. Route53 Zone
   ↓
8. Jenkins Resources (Key Pair, Security Group)
   ↓
9. VPC (Subnets, NAT, IGW)
```

---

## คำสั่งสำหรับการลบทั้งหมด (Copy-Paste)

```bash
# 1. ลบ ArgoCD Applications
kubectl delete application --all -n argocd

# 2. รอสักครู่ (2-3 นาที)
sleep 180

# 3. Destroy ทั้งหมด
cd Infrastructure/terraform
terraform destroy -auto-approve
```

---

## หมายเหตุ

- ใช้เวลาประมาณ **15-30 นาที** ในการลบทั้งหมด
- ตรวจสอบ AWS billing หลังจากลบเพื่อยืนยันว่าไม่มี resources ค้างอยู่
- เก็บ `terraform.tfstate` ไว้สำรองก่อนลบ (ถ้าต้องการกู้คืน)

---

## การสำรองข้อมูลก่อนลบ

```bash
# สำรอง terraform state
cp terraform.tfstate terraform.tfstate.backup.$(date +%Y%m%d_%H%M%S)

# สำรอง kubeconfig
kubectl config view --raw > kubeconfig.backup.$(date +%Y%m%d_%H%M%S).yaml

# Export ArgoCD applications
kubectl get applications -n argocd -o yaml > argocd-apps.backup.yaml
```

---

**สร้างโดย**: Medee DevOps Team  
**อัพเดทล่าสุด**: 2024

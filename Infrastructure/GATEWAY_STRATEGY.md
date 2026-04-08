# Gateway API Management Strategy

## ปัญหา: ทำไมต้องสร้าง Gateway API จาก Terraform?

### คำตอบสั้น
เพราะ **Gateway API CRDs** และ **AWS Load Balancer Controller** ต้องการ **IAM permissions** และ **infrastructure dependencies** ที่ ArgoCD จัดการไม่ได้

---

## Architecture Layers

```
┌─────────────────────────────────────────────────────────┐
│  Layer 1: Infrastructure (Terraform)                    │
│  - VPC, Subnets, NAT Gateway                           │
│  - EKS Cluster, Node Groups                            │
│  - IAM Roles, IRSA                                     │
│  - Security Groups                                     │
└─────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────┐
│  Layer 2: Platform Controllers (Terraform)              │
│  - Gateway API CRDs                                    │
│  - AWS Load Balancer Controller + IAM                  │
│  - Cert Manager                                        │
│  - ArgoCD                                              │
└─────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────┐
│  Layer 3: Platform Resources (ArgoCD)                   │
│  - Gateway instances                                   │
│  - GatewayClass configurations                         │
│  - ReferenceGrants                                     │
└─────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────┐
│  Layer 4: Applications (ArgoCD)                         │
│  - HTTPRoutes                                          │
│  - Services                                            │
│  - Deployments                                         │
└─────────────────────────────────────────────────────────┘
```

---

## ทำไม Gateway API ไม่สามารถอยู่ใน ArgoCD ทั้งหมด?

### 1. CRDs ต้องติดตั้งก่อน

```bash
# ถ้า ArgoCD พยายามสร้าง HTTPRoute แต่ CRDs ยังไม่มี
Error: no matches for kind "HTTPRoute" in version "gateway.networking.k8s.io/v1"
```

### 2. AWS Load Balancer Controller ต้องการ IAM Role

```yaml
# AWS LB Controller ต้องการ permissions เหล่านี้:
- elasticloadbalancing:*
- ec2:DescribeSubnets
- ec2:DescribeSecurityGroups
- iam:CreateServiceLinkedRole
```

ArgoCD ไม่สามารถสร้าง IAM Role ได้!

### 3. Chicken-and-Egg Problem

```
ArgoCD ต้องการ Gateway → Gateway ต้องการ LB Controller 
→ LB Controller ต้องการ IAM Role → IAM Role ต้องการ EKS
```

---

## แนวทางที่แนะนำ: Hybrid Approach

### Terraform จัดการ (Infrastructure + Controllers):

```hcl
# modules/gateway/main.tf

# 1. Install Gateway API CRDs
resource "helm_release" "gateway_api_crds" {
  name       = "gateway-api-crds"
  repository = "https://gateway-api.github.io/gateway-api"
  chart      = "gateway-api"
  namespace  = "gateway-system"
  create_namespace = true
}

# 2. Install AWS Load Balancer Controller with IAM
resource "helm_release" "aws_load_balancer_controller" {
  name       = "aws-load-balancer-controller"
  repository = "https://aws.github.io/eks-charts"
  chart      = "aws-load-balancer-controller"
  namespace  = "kube-system"

  set {
    name  = "serviceAccount.annotations.eks\\.amazonaws\\.com/role-arn"
    value = var.lb_controller_role_arn  # จาก EKS module
  }
}
```

### ArgoCD จัดการ (Gateway Instances + Applications):

```yaml
# gitops/infrastructure/gateway/gateway.yaml
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: main-gateway
  namespace: gateway-system
spec:
  gatewayClassName: aws-load-balancer-controller
  listeners:
  - name: http
    port: 80
    protocol: HTTP
  - name: https
    port: 443
    protocol: HTTPS
    tls:
      certificateRefs:
      - name: tls-cert
```

```yaml
# gitops/apps/medee-dev/httproute.yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: frontend-route
  namespace: medee-dev
spec:
  parentRefs:
  - name: main-gateway
    namespace: gateway-system
  hostnames:
  - "app.dodee.me"
  rules:
  - backendRefs:
    - name: frontend-service
      port: 8080
```

---

## ข้อดีของ Hybrid Approach

| Aspect | Terraform | ArgoCD |
|--------|-----------|--------|
| **CRDs** | ✅ ติดตั้ง | ❌ ใช้งาน |
| **Controllers** | ✅ Deploy + IAM | ❌ |
| **Gateway Instances** | ❌ | ✅ จัดการ |
| **HTTPRoutes** | ❌ | ✅ จัดการ |
| **Applications** | ❌ | ✅ จัดการ |

---

## ถ้าอยากให้ ArgoCD จัดการทั้งหมด (ไม่แนะนำ)

### ต้องทำแบบนี้:

1. **Terraform สร้าง EKS + IAM เท่านั้น**
2. **สร้าง ArgoCD App-of-Apps แบบ 2 ชั้น:**

```yaml
# bootstrap/infrastructure.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: infrastructure
  namespace: argocd
spec:
  project: default
  source:
    repoURL: git@github.com:Nut06/medee-gitops.git
    path: infrastructure
  destination:
    server: https://kubernetes.default.svc
    namespace: argocd
  syncPolicy:
    syncOptions:
    - CreateNamespace=true
    syncWaves:
    - wave: 0  # CRDs first
```

```yaml
# infrastructure/gateway-crds.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: gateway-crds
  namespace: argocd
  annotations:
    argocd.argoproj.io/sync-wave: "0"  # ติดตั้งก่อน
spec:
  source:
    repoURL: https://gateway-api.github.io/gateway-api
    chart: gateway-api
```

```yaml
# infrastructure/aws-lb-controller.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: aws-lb-controller
  namespace: argocd
  annotations:
    argocd.argoproj.io/sync-wave: "1"  # ติดตั้งหลัง CRDs
spec:
  source:
    repoURL: https://aws.github.io/eks-charts
    chart: aws-load-balancer-controller
  helm:
    parameters:
    - name: serviceAccount.annotations.eks\.amazonaws\.com/role-arn
      value: "arn:aws:iam::652253416828:role/medee-lb-controller"  # ต้องสร้างใน Terraform ก่อน
```

### ปัญหาของวิธีนี้:
- ❌ IAM Role ยังต้องสร้างใน Terraform อยู่ดี
- ❌ Sync waves ซับซ้อน
- ❌ ถ้า CRDs fail, ทุกอย่าง fail
- ❌ ยากต่อการ debug

---

## สรุป: ทำไมต้องแยก?

### Terraform (Infrastructure as Code):
- จัดการ **infrastructure** ที่ต้องการ AWS permissions
- สร้าง **foundation** ที่ stable
- ง่ายต่อการ destroy และสร้างใหม่

### ArgoCD (GitOps):
- จัดการ **applications** และ **configurations**
- Auto-sync จาก Git
- ง่ายต่อการ rollback

---

## Best Practice

```
Terraform:
├── VPC, EKS, IAM
├── Gateway API CRDs
├── AWS LB Controller
├── Cert Manager
└── ArgoCD

ArgoCD:
├── Gateway instances
├── HTTPRoutes
├── Applications
└── Configurations
```

**หลักการ**: ถ้า resource ต้องการ **AWS IAM** หรือเป็น **CRDs** → ใช้ **Terraform**  
ถ้าเป็น **application resources** → ใช้ **ArgoCD**

---

## คำแนะนำสำหรับ Medee Project

**ปัจจุบัน (ถูกต้องแล้ว)**:
```
Terraform → Gateway API CRDs + AWS LB Controller
ArgoCD → Applications + HTTPRoutes
```

**ถ้าอยากลด Terraform**:
```
Terraform → เฉพาะ VPC + EKS + IAM
ArgoCD → ทุกอย่างที่เหลือ (แต่ซับซ้อนกว่า)
```

**แนะนำ**: ใช้แบบปัจจุบันต่อไป เพราะ:
- ✅ Clear separation of concerns
- ✅ ง่ายต่อการ debug
- ✅ Stable infrastructure
- ✅ Flexible application deployment

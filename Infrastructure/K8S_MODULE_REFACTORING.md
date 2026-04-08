# K8s Module Refactoring Guide

## ปัญหา: อะไรควรอยู่ใน Terraform vs GitOps?

---

## แนวทาง 1: Minimal Terraform (แนะนำ)

### Terraform เก็บเฉพาะ Bootstrap Resources

```hcl
# modules/k8s/main.tf (Minimal Version)

locals {
  common_tags = {
    Environment = var.environment
    Project     = var.project_name
    ManagedBy   = "terraform"
    Component   = "argocd"
  }
}

# -----------------------------------------------
# 1. Infrastructure Namespaces (จำเป็น)
# -----------------------------------------------
resource "kubernetes_namespace_v1" "infrastructure_namespaces" {
  for_each = toset([
    "argocd",           # ArgoCD ต้องมีก่อน
    "gateway-system"    # Gateway Controller ต้องมีก่อน
  ])
  
  metadata {
    name = each.value
    labels = {
      name = each.value
      "app.kubernetes.io/managed-by" = "terraform"
      environment = var.environment
      project     = var.project_name
      purpose     = "infrastructure"
    }
  }
}

# -----------------------------------------------
# 2. ArgoCD (Bootstrap Tool)
# -----------------------------------------------
resource "helm_release" "argocd" {
  name       = "argocd"
  repository = "https://argoproj.github.io/argo-helm"
  chart      = "argo-cd"
  version    = "9.4.16"
  namespace  = "argocd"

  wait            = true
  wait_for_jobs   = true
  timeout         = 600
  atomic          = true
  cleanup_on_fail = true

  values = [
    yamlencode({
      server = {
        service = {
          type = "ClusterIP"
        }
        ingress = {
          enabled = false
        }
      }
    })
  ]

  depends_on = [kubernetes_namespace_v1.infrastructure_namespaces]
}

# -----------------------------------------------
# 3. Root Application (Bootstrap GitOps)
# -----------------------------------------------
resource "kubernetes_manifest" "root_application" {
  manifest = {
    apiVersion = "argoproj.io/v1alpha1"
    kind       = "Application"
    metadata = {
      name      = "root-app"
      namespace = "argocd"
      finalizers = [
        "resources-finalizer.argocd.argoproj.io"
      ]
    }
    spec = {
      project = "default"
      source = {
        repoURL        = var.gitops_repo_url
        targetRevision = var.gitops_target_revision
        path           = "bootstrap"
      }
      destination = {
        server    = "https://kubernetes.default.svc"
        namespace = "argocd"
      }
      syncPolicy = {
        automated = {
          prune    = true
          selfHeal = true
        }
        syncOptions = [
          "CreateNamespace=true",
          "ServerSideApply=true"
        ]
      }
    }
  }

  depends_on = [helm_release.argocd]
}

# -----------------------------------------------
# 4. ArgoCD RBAC (จำเป็นสำหรับ cluster-admin)
# -----------------------------------------------
resource "kubernetes_cluster_role_v1" "argocd_application_controller" {
  metadata {
    name = "argocd-application-controller-cluster-admin"
    labels = merge(local.common_tags, {
      "app.kubernetes.io/name"      = "argocd"
      "app.kubernetes.io/component" = "application-controller"
    })
  }

  rule {
    api_groups = ["*"]
    resources  = ["*"]
    verbs      = ["*"]
  }
}

resource "kubernetes_cluster_role_binding_v1" "argocd_application_controller" {
  metadata {
    name = "argocd-application-controller-cluster-admin"
    labels = merge(local.common_tags, {
      "app.kubernetes.io/name"      = "argocd"
      "app.kubernetes.io/component" = "application-controller"
    })
  }

  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = kubernetes_cluster_role_v1.argocd_application_controller.metadata[0].name
  }

  subject {
    kind      = "ServiceAccount"
    name      = "argocd-application-controller"
    namespace = "argocd"
  }

  depends_on = [helm_release.argocd]
}
```

### GitOps จัดการที่เหลือ

```yaml
# gitops/infrastructure/namespaces/namespaces.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: medee-dev
  labels:
    name: medee-dev
    environment: dev
    project: medee
    purpose: application
---
apiVersion: v1
kind: Namespace
metadata:
  name: medee-staging
  labels:
    name: medee-staging
    environment: staging
    project: medee
    purpose: application
---
apiVersion: v1
kind: Namespace
metadata:
  name: monitoring
  labels:
    name: monitoring
    environment: dev
    project: medee
    purpose: infrastructure
---
apiVersion: v1
kind: Namespace
metadata:
  name: vault
  labels:
    name: vault
    environment: dev
    project: medee
    purpose: infrastructure
```

```yaml
# gitops/infrastructure/reference-grants/reference-grants.yaml
apiVersion: gateway.networking.k8s.io/v1beta1
kind: ReferenceGrant
metadata:
  name: gateway-reference-grant
  namespace: medee-dev
spec:
  from:
  - group: gateway.networking.k8s.io
    kind: HTTPRoute
    namespace: medee-dev
  to:
  - group: gateway.networking.k8s.io
    kind: Gateway
---
apiVersion: gateway.networking.k8s.io/v1beta1
kind: ReferenceGrant
metadata:
  name: gateway-reference-grant
  namespace: medee-staging
spec:
  from:
  - group: gateway.networking.k8s.io
    kind: HTTPRoute
    namespace: medee-staging
  to:
  - group: gateway.networking.k8s.io
    kind: Gateway
---
apiVersion: gateway.networking.k8s.io/v1beta1
kind: ReferenceGrant
metadata:
  name: gateway-reference-grant
  namespace: argocd
spec:
  from:
  - group: gateway.networking.k8s.io
    kind: HTTPRoute
    namespace: argocd
  to:
  - group: gateway.networking.k8s.io
    kind: Gateway
---
apiVersion: gateway.networking.k8s.io/v1beta1
kind: ReferenceGrant
metadata:
  name: gateway-reference-grant
  namespace: monitoring
spec:
  from:
  - group: gateway.networking.k8s.io
    kind: HTTPRoute
    namespace: monitoring
  to:
  - group: gateway.networking.k8s.io
    kind: Gateway
```

```yaml
# gitops/infrastructure/vault/vault-app.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: vault
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://helm.releases.hashicorp.com
    chart: vault
    targetRevision: 0.28.0
    helm:
      values: |
        server:
          ha:
            enabled: false
        injector:
          enabled: true
  destination:
    server: https://kubernetes.default.svc
    namespace: vault
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
    - CreateNamespace=true
```

---

## แนวทาง 2: Hybrid (ถ้าต้องการความปลอดภัยมากขึ้น)

### Terraform จัดการ Critical Resources

```hcl
# modules/k8s/main.tf (Hybrid Version)

# 1. ทุก Namespaces (รวม application namespaces)
resource "kubernetes_namespace_v1" "all_namespaces" {
  for_each = toset([
    "argocd",
    "gateway-system",
    "medee-dev",
    "medee-staging",
    "monitoring",
    "vault"
  ])
  
  metadata {
    name = each.value
    labels = {
      name = each.value
      environment = var.environment
      project = var.project_name
    }
  }
}

# 2. ArgoCD + Root App
# (เหมือนแนวทาง 1)

# 3. Vault (ถ้าเป็น critical infrastructure)
resource "helm_release" "vault" {
  name       = "vault"
  repository = "https://helm.releases.hashicorp.com"
  chart      = "vault"
  namespace  = "vault"

  values = [yamlencode({
    server = {
      ha = { enabled = false }
    }
    injector = {
      enabled = true
    }
  })]

  depends_on = [kubernetes_namespace_v1.all_namespaces]
}
```

### GitOps จัดการ Application Resources

```yaml
# gitops/infrastructure/reference-grants/
# (เหมือนแนวทาง 1)
```

---

## เปรียบเทียบแนวทาง

| Resource | แนวทาง 1 (Minimal) | แนวทาง 2 (Hybrid) |
|----------|-------------------|------------------|
| **argocd namespace** | Terraform ✅ | Terraform ✅ |
| **gateway-system namespace** | Terraform ✅ | Terraform ✅ |
| **app namespaces** | GitOps ✅ | Terraform ✅ |
| **ArgoCD Helm** | Terraform ✅ | Terraform ✅ |
| **root-app** | Terraform ✅ | Terraform ✅ |
| **ReferenceGrants** | GitOps ✅ | GitOps ✅ |
| **Vault** | GitOps ✅ | Terraform ✅ |
| **RBAC** | Terraform ✅ | Terraform ✅ |

---

## ข้อดี-ข้อเสีย

### แนวทาง 1: Minimal Terraform

**ข้อดี**:
- ✅ Terraform น้อยที่สุด
- ✅ GitOps จัดการทุกอย่างที่เป็น application
- ✅ ง่ายต่อการแก้ไข (แก้ใน Git)
- ✅ True GitOps philosophy

**ข้อเสีย**:
- ⚠️ Namespaces อาจถูกลบโดยไม่ตั้งใจ
- ⚠️ ต้องรอ ArgoCD sync ก่อน

### แนวทาง 2: Hybrid

**ข้อดี**:
- ✅ Namespaces stable (ไม่ถูกลบง่าย)
- ✅ Critical infrastructure ใน Terraform
- ✅ ปลอดภัยกว่า

**ข้อเสีย**:
- ❌ Terraform มากขึ้น
- ❌ แก้ไขต้องรัน terraform apply

---

## คำแนะนำสำหรับ Medee Project

### ใช้ แนวทาง 1: Minimal Terraform ✅

**เหตุผล**:
1. คุณใช้ GitOps อยู่แล้ว
2. Namespaces ไม่ใช่ critical resource
3. ง่ายต่อการจัดการและแก้ไข
4. ลด Terraform state complexity

### Migration Steps:

1. **ย้าย Namespaces ไป GitOps**
   ```bash
   # ลบจาก Terraform state (ไม่ลบ resource จริง)
   terraform state rm 'kubernetes_namespace_v1.infrastructure_namespaces["medee-dev"]'
   terraform state rm 'kubernetes_namespace_v1.infrastructure_namespaces["medee-staging"]'
   terraform state rm 'kubernetes_namespace_v1.infrastructure_namespaces["monitoring"]'
   terraform state rm 'kubernetes_namespace_v1.infrastructure_namespaces["vault"]'
   ```

2. **ย้าย ReferenceGrants ไป GitOps**
   ```bash
   terraform state rm 'kubernetes_manifest.reference_grants["medee-dev"]'
   terraform state rm 'kubernetes_manifest.reference_grants["medee-staging"]'
   terraform state rm 'kubernetes_manifest.reference_grants["argocd"]'
   terraform state rm 'kubernetes_manifest.reference_grants["monitoring"]'
   ```

3. **ย้าย Vault ไป GitOps**
   ```bash
   terraform state rm 'helm_release.vault'
   ```

4. **สร้าง GitOps manifests**
   - สร้างไฟล์ใน `medee-gitops/infrastructure/`

5. **Apply ใน Terraform**
   ```bash
   terraform apply
   ```

---

## สรุป

**ใช่ครับ! ส่วนใหญ่ใน k8s module ย้ายไป GitOps ได้**

**ควรเก็บใน Terraform**:
- ArgoCD Helm Release
- root-app
- argocd namespace
- gateway-system namespace
- ArgoCD RBAC

**ควรย้ายไป GitOps**:
- Application namespaces (medee-dev, medee-staging, monitoring, vault)
- ReferenceGrants
- Vault
- Application-level RBAC

**ผลลัพธ์**: Terraform เบาลง, GitOps จัดการ application layer ทั้งหมด ✅

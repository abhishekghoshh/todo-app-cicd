# Getting Started

This guide will walk you through setting up and deploying the Todo App CI/CD pipeline.

## Quick Start (Recommended)

### 1. Clone the Repository
```bash
git clone https://github.com/abhishekghoshh/todo-app-cicd
cd todo-app-cicd
```

### 2. Deploy to Kubernetes with One Command
```bash
# Deploy to feature environment (or dev, staging, prod)
./argo-deploy.sh feature
```

**What this script does automatically:**
- ✅ Installs ArgoCD if not present
- ✅ Configures Helm support in ArgoCD
- ✅ Deploys the root application pointing to your environment
- ✅ Sets up port-forwarding to ArgoCD UI
- ✅ Provides admin credentials for immediate access

### 3. Access Your Application

```bash
# The script automatically starts port-forwarding
# ArgoCD UI: http://localhost:8080 (admin/password from script output)

# Access your Todo app (after ArgoCD syncs):
kubectl port-forward svc/todo-app-svc -n todo-app 3000:8080
# Todo App: http://localhost:3000
```

## Local Development Setup

### Using Docker Compose

```bash
# Run locally with Docker Compose
cp .env.example .env  # Configure environment variables
docker-compose up -d

# Access at http://localhost:8080
# Login: admin / admin-password
```

### Default Credentials
- **Username**: `admin`
- **Password**: `admin-password`

## Manual Deployment

### Step 1: Build and Push Images

```bash
# Build multi-architecture Docker image
docker buildx build --platform linux/amd64,linux/arm64 \
  -t abhishek1009/todo-app:latest --push .

# Package and push Helm chart to OCI registry
helm package helm --destination ./packages
helm push ./packages/todo-app-chart-1.0.0.tgz oci://registry-1.docker.io/abhishek1009
```

### Step 2: Install ArgoCD

```bash
# Install ArgoCD
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Wait for ArgoCD to be ready
kubectl wait deployment/argocd-server \
  -n argocd \
  --for=condition=Available=True \
  --timeout=180s
```

### Step 3: Enable Helm Support

```bash
./enable-helm.sh
```

### Step 4: Apply Root Application

```bash
# Set your target environment
export targetEnvironment=feature

# Apply the root application
envsubst '${targetEnvironment}' < application.yaml | kubectl apply -f -
```

### Step 5: Access ArgoCD UI

```bash
# Port forward to access UI
kubectl port-forward svc/argocd-server -n argocd 8080:443

# Get initial admin password
kubectl -n argocd get secret argocd-initial-admin-secret \
  -o jsonpath="{.data.password}" | base64 -d
```

## Verify Deployment

```bash
# Check ArgoCD applications
kubectl get applications -n argocd

# Check application pods
kubectl get pods -n todo-app

# View application logs
kubectl logs -f deployment/todo-app -n todo-app
```

## Next Steps

- [Configure environments](configuration.md)
- [Set up CI/CD workflows](cicd-workflows.md)
- [Monitor your deployment](monitoring.md)

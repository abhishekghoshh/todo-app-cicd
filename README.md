# Todo App CI/CD Pipeline

A comprehensive **GitOps-powered CI/CD pipeline** for a modern Go-based Todo application with complete Kubernetes deployment automation using ArgoCD, Helm, and GitHub Actions.

## 🚀 What This Project Achieves

### **Complete End-to-End CI/CD Pipeline**
This project demonstrates a **production-ready CI/CD pipeline** that automatically builds, tests, packages, and deploys a Go web application to Kubernetes using modern cloud-native practices.

### **Technology Stack**
- **Backend**: Go 1.25+ with native HTTP routing, MongoDB integration
- **Frontend**: HTMX-powered interactive UI with server-side rendering
- **Database**: MongoDB with session management and TTL indexes
- **Containerization**: Multi-stage Docker builds with Alpine Linux
- **Orchestration**: Kubernetes with Helm charts and Kustomize overlays
- **GitOps**: ArgoCD for automated deployment and management
- **CI/CD**: GitHub Actions for build automation and artifact publishing

### **Key Features Implemented**
- ✅ **Multi-environment deployment** (dev, feature, staging, production)
- ✅ **Automated Docker image builds** with multi-architecture support (amd64/arm64)
- ✅ **Helm chart packaging** and publishing to OCI registry (Docker Hub)
- ✅ **GitOps workflow** with ArgoCD for deployment automation
- ✅ **Environment-specific configurations** using Kustomize overlays
- ✅ **Persistent storage** for user uploads and database data
- ✅ **Ingress configuration** with load balancing and SSL termination support
- ✅ **Resource management** with CPU/memory limits and requests
- ✅ **Session-based authentication** with secure cookie handling
- ✅ **File upload functionality** with persistent volume storage

## 🏗️ Architecture Overview

### **Application Architecture**
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTMX Client   │────│   Go Backend    │────│    MongoDB      │
│  (Dashboard)    │    │  (REST + Auth)  │    │  (Persistence)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### **CI/CD Pipeline Flow**
```
┌─────────────┐   ┌──────────────┐   ┌─────────────┐   ┌─────────────┐
│ Code Push   │──▶│ GitHub       │──▶│ ArgoCD      │──▶│ Kubernetes  │
│ (feature)   │   │ Actions      │   │ Sync        │   │ Deployment  │
│             │   │ (Build/Push) │   │ (GitOps)    │   │ (Live App)  │
└─────────────┘   └──────────────┘   └─────────────┘   └─────────────┘
```

## 📁 Project Structure Deep Dive

### **Root Directory Management Scripts**

#### **`application.yaml`** - ArgoCD Root Application
```yaml
# Root application that manages all environment deployments
# Uses environment variable substitution for dynamic targeting
```
This file defines the **"App of Apps" pattern** in ArgoCD, allowing you to manage multiple applications across different environments from a single root application. It uses environment variable substitution (`${targetEnvironment}`) to dynamically target different branches and paths.

#### **`argo-deploy.sh`** - Automated ArgoCD Deployment Script
```bash
# Usage: ./argo-deploy.sh <environment>
# Examples: 
#   ./argo-deploy.sh feature    # Deploy feature environment
#   ./argo-deploy.sh dev        # Deploy dev environment  
#   ./argo-deploy.sh destroy    # Remove ArgoCD completely
```
**Key capabilities:**
- ✅ **One-command deployment** to any environment
- ✅ **ArgoCD installation** and configuration automation
- ✅ **Helm integration** setup within ArgoCD
- ✅ **Environment variable substitution** for dynamic configuration
- ✅ **Port-forwarding** automation for immediate UI access
- ✅ **Complete cleanup** functionality with destroy option

#### **`enable-helm.sh`** - ArgoCD Helm Integration Setup
```bash
# Automatically configures ArgoCD to support Helm within Kustomize
# Patches argocd-cm ConfigMap with necessary build options
```
**What it does:**
- ✅ **Enables Helm support** in Kustomize for ArgoCD
- ✅ **Configures build options** (`--enable-helm --load-restrictor LoadRestrictionsNone`)
- ✅ **Restarts ArgoCD components** to apply changes
- ✅ **Waits for readiness** before completion

### **Application Structure**
```
├── src/main/                    # Go application source code
│   ├── main.go                 # Application entry point with routing
│   ├── handlers.go             # HTTP handlers and middleware
│   ├── models.go               # Data models (User, Todo, Session)
│   └── db.go                   # MongoDB operations and connection
├── src/resources/               # Frontend templates and assets
│   └── templates/              # HTMX-powered HTML templates
├── helm/                       # Production Helm charts
│   ├── Chart.yaml              # Chart metadata
│   ├── values.yaml             # Default configuration values
│   └── templates/              # Kubernetes manifest templates
├── argocd/                     # GitOps configuration
│   ├── apps/todo-app/          # Application-specific configs
│   │   ├── base/               # Base Kustomize configuration
│   │   └── overlays/           # Environment-specific overrides
│   └── argo-apps/              # ArgoCD Application definitions
├── .github/workflows/          # CI/CD automation
│   ├── docker.yaml             # Docker image build and push
│   └── helm-chart.yaml         # Helm chart packaging and publish
```

## 🛠️ Prerequisites
- Docker installed and running
- Kubernetes cluster (local or cloud)
- kubectl configured and connected to cluster
- Helm 3.x installed
- Git repository access
- Docker Hub account (for image and chart registry)

## 🚀 Quick Start Guide

### **Option 1: Automated Deployment (Recommended)**

#### 1. Clone and Setup
```bash
git clone https://github.com/abhishekghoshh/todo-app-cicd
cd todo-app-cicd
```

#### 2. Deploy to Kubernetes with One Command
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

#### 3. Access Your Application
```bash
# The script automatically starts port-forwarding
# ArgoCD UI: http://localhost:8080 (admin/password from script output)

# Access your Todo app (after ArgoCD syncs):
kubectl port-forward svc/todo-app-svc -n todo-app 3000:8080
# Todo App: http://localhost:3000
```

### **Option 2: Manual Step-by-Step Setup**

#### 1. Local Development
```bash
# Run locally with Docker Compose
cp .env.example .env  # Configure environment variables
docker-compose up -d

# Access at http://localhost:8080
# Login: admin / admin-password
```

#### 2. Build and Push Images (CI/CD handles this automatically)
```bash
# Build multi-architecture Docker image
docker buildx build --platform linux/amd64,linux/arm64 \
  -t abhishek1009/todo-app:latest --push .

# Package and push Helm chart to OCI registry
helm package helm --destination ./packages
helm push ./packages/todo-app-chart-1.0.0.tgz oci://registry-1.docker.io/abhishek1009
```

#### 3. Manual ArgoCD Setup
```bash
# Install ArgoCD
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Enable Helm support
./enable-helm.sh

# Apply root application
export targetEnvironment=feature
envsubst '${targetEnvironment}' < application.yaml | kubectl apply -f -
```

## 🔧 Advanced Configuration

### **Environment Management**

Each environment has its own configuration overlay:

```bash
# Feature environment
argocd/apps/todo-app/overlays/feature/values.yaml

# Development environment  
argocd/apps/todo-app/overlays/dev/values.yaml
```

**Key configuration options:**
- Image tags and repositories
- Resource limits and requests
- Ingress hostnames and TLS settings
- Storage class and volume sizes
- Environment variables

### **GitHub Actions CI/CD Automation**

#### **🐳 Docker Image Pipeline** (`.github/workflows/docker.yaml`)
**Purpose**: Automatically builds and publishes container images whenever application code changes.

**Trigger Conditions**:
- Push to `feature` branch with changes in: `src/**`, `Dockerfile`, `go.mod`, `go.sum`
- Manual workflow dispatch available

**What it does**:
- ✅ **Multi-stage Docker build** using optimized Go builder pattern
- ✅ **Multi-architecture support** (linux/amd64, linux/arm64) for broad compatibility
- ✅ **Docker Hub authentication** using GitHub secrets
- ✅ **Image tagging strategy**: 
  - `latest` tag for the most recent build
  - `${GITHUB_SHA}` tag for specific commit traceability
- ✅ **Build caching** to speed up subsequent builds
- ✅ **Automated push** to Docker Hub registry (`abhishek1009/todo-app`)

**Key Features**:
```yaml
# Builds for multiple platforms simultaneously
platforms: linux/amd64,linux/arm64

# Dual tagging for flexibility
tags: |
  abhishek1009/todo-app:latest
  abhishek1009/todo-app:${GITHUB_SHA}

# Registry caching for faster builds
cache-from: type=registry,ref=abhishek1009/todo-app:cache
```

#### **📦 Helm Chart Pipeline** (`.github/workflows/helm-chart.yaml`)
**Purpose**: Packages and publishes Helm charts to OCI registry whenever chart configurations change.

**Trigger Conditions**:
- Push to `feature` branch with changes in: `helm/**` directory
- Manual workflow dispatch available

**What it does**:
- ✅ **Helm chart packaging** using chart metadata from `Chart.yaml`
- ✅ **OCI registry publishing** to Docker Hub (modern Helm 3 approach)
- ✅ **Dual versioning strategy**:
  - **Semantic version** from `Chart.yaml` (e.g., `1.0.0`)
  - **Commit-based version** for development tracking (e.g., `1.0.0-a1b2c3d`)
- ✅ **Automated chart validation** and packaging
- ✅ **Docker Hub OCI registry** storage for charts

**Chart Publishing Process**:
```bash
# Extract chart metadata
CHART_NAME=$(grep 'name:' helm/Chart.yaml | awk '{print $2}')
CHART_VERSION=$(grep 'version:' helm/Chart.yaml | awk '{print $2}')

# Package and push stable version
helm package helm --destination ./packages
helm push ./packages/$CHART_NAME-$CHART_VERSION.tgz oci://registry-1.docker.io/abhishek1009

# Create development version with commit SHA
CHART_VERSION_FROM_COMMIT="${CHART_VERSION}-${GITHUB_SHA::7}"
helm package helm --version $CHART_VERSION_FROM_COMMIT --destination ./packages
helm push ./packages/$CHART_NAME-$CHART_VERSION_FROM_COMMIT.tgz oci://registry-1.docker.io/abhishek1009
```

#### **🔄 Complete CI/CD Flow**
```
Code Change → GitHub Push → Actions Triggered
     ↓
┌─────────────────────┬─────────────────────┐
│   docker.yaml       │   helm-chart.yaml   │
│   ===============   │   =================  │
│ • Build Go app      │ • Package Helm chart│
│ • Create container  │ • Extract versions   │
│ • Multi-arch build  │ • Push to OCI registry│
│ • Push to registry  │ • Create dev version │
└─────────────────────┴─────────────────────┘
     ↓
ArgoCD Detects Changes → Pulls New Artifacts → Deploys to Kubernetes
```

**Why This Approach Works**:
- 🎯 **Separation of Concerns**: Different pipelines for different artifact types
- 🚀 **Fast Feedback**: Only builds what changed (source vs config)
- 🏗️ **Artifact Traceability**: Every build tagged with commit SHA
- 🔒 **Secure Automation**: Uses GitHub secrets for registry authentication
- 🌐 **Multi-Platform Support**: Images work on different architectures
- 📦 **Modern Standards**: Uses OCI registries for both images and charts

### **Monitoring and Troubleshooting**

```bash
# Check ArgoCD application status
kubectl get applications -n argocd

# Monitor application pods
kubectl get pods -n todo-app -w

# View application logs
kubectl logs -f deployment/todo-app -n todo-app

# Check persistent volumes
kubectl get pv,pvc -n todo-app

# ArgoCD CLI operations
argocd app list
argocd app sync todo-app
argocd app history todo-app
```

### **Cleanup and Destroy**

```bash
# Remove specific environment
kubectl delete application todo-app -n argocd

# Complete ArgoCD removal
./argo-deploy.sh destroy
```

## 🏆 Production Deployment Features

### **High Availability Setup**
- **Persistent storage** for user uploads and MongoDB data
- **Resource management** with CPU/memory limits
- **Health checks** and readiness probes
- **Horizontal scaling** capabilities via Kubernetes HPA

### **Security Implementations**
- **HTTP-only cookies** for session management
- **Session TTL** with automatic cleanup
- **Container security** with non-root user and minimal Alpine image
- **Network policies** support (configurable)
- **TLS/SSL** ready ingress configuration

### **DevOps Best Practices**
- **GitOps methodology** with ArgoCD
- **Infrastructure as Code** with Helm and Kustomize
- **Multi-environment promotion** pipeline
- **Automated rollback** capabilities
- **Declarative configuration** management
- **OCI registry** for both images and Helm charts

## 🎯 What You Learn From This Project

### **Modern Cloud-Native Patterns**
1. **GitOps Workflow**: Code push → Automated build → GitOps sync → Kubernetes deployment
2. **App of Apps Pattern**: Managing multiple applications through a single root ArgoCD application
3. **Kustomize + Helm Hybrid**: Combining the power of both tools for maximum flexibility
4. **OCI Registry Usage**: Storing both Docker images and Helm charts in the same registry

### **Production-Ready Practices**
1. **Multi-stage Docker builds** for optimized container images
2. **Environment-specific configurations** without code duplication
3. **Persistent storage strategies** for stateful applications
4. **Resource management** and horizontal scaling
5. **Security hardening** with proper authentication and session management

### **Automation Excellence**
1. **One-command deployment** to any environment
2. **Automated CI/CD pipelines** triggered by code changes
3. **Self-healing deployments** with ArgoCD sync policies
4. **Infrastructure automation** with shell scripting and environment substitution

This project serves as a **comprehensive template** for building production-grade CI/CD pipelines for modern web applications, demonstrating enterprise-level practices in a practical, hands-on implementation.


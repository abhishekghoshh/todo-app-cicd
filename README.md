# Todo App CI/CD Pipeline

A comprehensive **GitOps-powered CI/CD pipeline** for a modern Go-based Todo application with complete Kubernetes deployment automation using ArgoCD, Helm, and GitHub Actions.

## Referrence

- [Simplifying Argo CD Deployments with Kustomize: 10 Benefits and 3 AWS Scenarios](https://devopsarcheologicst.medium.com/simplifying-argo-cd-deployments-with-kustomize-10-benefits-and-3-aws-scenarios-ec6494ffd7fc)


## 📚 Documentation

Complete documentation is available in the `docs/` directory and can be viewed using MkDocs:

```bash
# Install MkDocs
pip install mkdocs-material

# Serve documentation locally
mkdocs serve

# Access at http://127.0.0.1:8000
```

### Quick Documentation Links

- **[Getting Started](docs/getting-started.md)** - Quick setup and deployment guide
- **[Project Structure](docs/project-structure.md)** - Directory layout and organization
- **[Deployment Scripts](docs/deployment-scripts.md)** - Automated deployment tools
- **[CI/CD Workflows](docs/cicd-workflows.md)** - GitHub Actions automation
- **[Configuration](docs/configuration.md)** - Environment management
- **[Monitoring](docs/monitoring.md)** - Troubleshooting and debugging
- **[Production Features](docs/production-features.md)** - Production capabilities
- **[Learning Outcomes](docs/learning-outcomes.md)** - What you'll learn

## 🚀 Quick Start

### One-Command Deployment

```bash
# Clone repository
git clone https://github.com/abhishekghoshh/todo-app-cicd
cd todo-app-cicd

# Deploy to Kubernetes
./argo-deploy.sh feature
```

**That's it!** The script will:
- ✅ Install ArgoCD if needed
- ✅ Configure Helm support
- ✅ Deploy your application
- ✅ Provide access credentials

### Access Your Application

```bash
# ArgoCD UI (auto-started by script)
http://localhost:8080

# Todo App (after ArgoCD syncs)
kubectl port-forward svc/todo-app-svc -n todo-app 3000:8080
http://localhost:3000
```

**Default credentials**: `admin / admin-password`

## 📖 Full Documentation

For detailed guides, configuration options, and troubleshooting, see the [complete documentation](docs/).

## 🎯 What This Project Achieves

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
│   ├── helm-chart.yaml         # Helm chart packaging and publish
│   └── sync-to-dev.yaml        # ArgoCD config sync automation
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

This project implements a **sophisticated three-workflow GitOps pipeline** that handles application builds, chart packaging, and configuration synchronization with intelligent pull request management.

#### **🐳 Workflow 1: Docker Image Pipeline** (`.github/workflows/docker.yaml`)
**Purpose**: Builds and publishes container images, then automatically updates environment configurations and creates PRs.

**Trigger Conditions**:
- Push to `feature` branch with changes in: `src/**`, `Dockerfile`, `go.mod`, `go.sum`
- Manual workflow dispatch available

**Complete Workflow Steps**:
1. **Build Phase**:
   - ✅ Multi-stage Docker build using optimized Go builder pattern
   - ✅ Multi-architecture support (linux/amd64, linux/arm64)
   - ✅ Docker Hub authentication using GitHub secrets
   - ✅ Image tagging: `latest` and `${GITHUB_SHA}` for traceability
   - ✅ Build caching for faster subsequent builds
   - ✅ Automated push to Docker Hub (`abhishek1009/todo-app`)

2. **Update Phase** (Critical GitOps Integration):
   - ✅ **Updates image tags** in `argocd/apps/todo-app/overlays/*/values.yaml` for dev, uat, and prod
   - ✅ **Commits using Docker Bot account** with personal access token
   - ✅ **Pushes changes** to feature branch
   
3. **Pull Request Phase**:
   - ✅ **Creates/updates PR** from `feature` → `dev` branch
   - ✅ **Includes image digest and commit SHA** in PR description
   - ✅ **Labels PR** with `automated-pr` and `docker-update`

**Why PR is necessary despite sync-to-dev workflow**:
```yaml
# Docker bot account pushes prevent infinite workflow loops
# GitHub Actions won't trigger workflows on commits from the bot account
# This prevents cyclic executions: docker.yaml → push → docker.yaml → ...
```

#### **📦 Workflow 2: Helm Chart Pipeline** (`.github/workflows/helm-chart.yaml`)
**Purpose**: Packages and publishes Helm charts to OCI registry, then updates Kustomize references and creates PRs.

**Trigger Conditions**:
- Push to `feature` branch with changes in: `helm/**` directory
- Manual workflow dispatch available

**Complete Workflow Steps**:
1. **Package Phase**:
   - ✅ Helm chart packaging from `helm/` directory
   - ✅ Version extraction from `Chart.yaml` (e.g., `1.0.0`)
   - ✅ **Dual versioning strategy**:
     - Semantic version for stable releases
     - Commit-based version for development (`1.0.0-a1b2c3d`)
   - ✅ OCI registry publishing to Docker Hub

2. **Update Phase** (Kustomize Integration):
   - ✅ **Updates Helm chart references** in `argocd/apps/todo-app/overlays/*/kustomization.yaml`
   - ✅ **Updates chart versions** for dev, uat, and prod environments
   - ✅ **Commits using Docker Bot account** with personal access token
   - ✅ **Pushes changes** to feature branch

3. **Pull Request Phase**:
   - ✅ **Creates/updates PR** from `feature` → `dev` branch
   - ✅ **Includes chart version and commit SHA** in PR description
   - ✅ **Labels PR** with `automated-pr` and `helm-update`

**Chart Publishing Process**:
```bash
# Extract metadata
CHART_NAME=$(grep 'name:' helm/Chart.yaml | awk '{print $2}')
CHART_VERSION=$(grep 'version:' helm/Chart.yaml | awk '{print $2}')

# Push stable version
helm push ./packages/$CHART_NAME-$CHART_VERSION.tgz oci://registry-1.docker.io/abhishek1009

# Push development version with commit SHA
CHART_VERSION_DEV="${CHART_VERSION}-${GITHUB_SHA::7}"
helm push ./packages/$CHART_NAME-$CHART_VERSION_DEV.tgz oci://registry-1.docker.io/abhishek1009
```

**Why PR is necessary despite redundancy**:
- Both workflows update files in `argocd/` directory
- Using Docker Bot account prevents infinite workflow loops
- GitHub won't trigger `sync-to-dev.yaml` on bot commits
- PR creation ensures human review before merging to dev

#### **🔄 Workflow 3: ArgoCD Sync Pipeline** (`.github/workflows/sync-to-dev.yaml`)
**Purpose**: Detects any changes in ArgoCD configurations and creates PRs to synchronize feature branch to dev.

**Trigger Conditions**:
- Push to `feature` branch with changes in: `argocd/**` directory
- Manual workflow dispatch available

**Complete Workflow Steps**:
1. **Detection Phase**:
   - ✅ Monitors all files under `argocd/` directory
   - ✅ Triggers on any configuration changes (values, kustomization, app definitions)
   
2. **Pull Request Phase**:
   - ✅ **Creates/updates PR** from `feature` → `dev` branch
   - ✅ **Includes changed files list** in PR description
   - ✅ **Labels PR** with `automated-pr` and `argocd-sync`
   - ✅ **Consolidates multiple changes** into single PR

**Why this workflow exists**:
- Catches manual configuration changes in `argocd/` directory
- Serves as backup for Docker/Helm workflows
- Ensures all ArgoCD changes go through PR review process
- Would normally catch Docker/Helm workflow updates, but prevented by bot account

#### **🎯 The Complete CI/CD Flow with Pull Request Strategy**

```
Code Change (src/ or helm/) → Push to Feature Branch
     ↓
┌─────────────────────┬─────────────────────┬─────────────────────┐
│   docker.yaml       │   helm-chart.yaml   │   sync-to-dev.yaml  │
│   ===============   │   =================  │   ================  │
│ 1. Build image      │ 1. Package chart    │ 1. Detect argocd/   │
│ 2. Push to Hub      │ 2. Push to OCI      │    changes          │
│ 3. Update values    │ 3. Update kustomize │ 2. Create/update PR │
│ 4. Commit as bot    │ 4. Commit as bot    │    feature→dev      │
│ 5. Push to feature  │ 5. Push to feature  │ 3. Request review   │
│ 6. Create PR        │ 6. Create PR        │                     │
│    feature→dev      │    feature→dev      │                     │
└─────────────────────┴─────────────────────┴─────────────────────┘
     ↓                      ↓                      ↓
     └──────────────────────┴──────────────────────┘
                            ↓
              Pull Request Created/Updated
              (feature → dev branch)
                            ↓
              Manual Review & Approval
                            ↓
              Merge to Dev Branch
                            ↓
              ArgoCD Detects Dev Branch Change
                            ↓
              Pulls Latest Artifacts
              (Docker images + Helm charts)
                            ↓
              Deploys to Kubernetes Cluster
```

#### **🔐 Preventing Infinite Workflow Loops**

**The Problem**:
```
docker.yaml runs → updates argocd/ files → commits & pushes
    ↓
sync-to-dev.yaml detects argocd/ change → creates PR
    ↓
PR merged → triggers docker.yaml again → INFINITE LOOP! ❌
```

**The Solution**:
```yaml
# All three workflows use Docker Bot account for commits
# GitHub Actions configuration prevents bot-triggered workflows

# In docker.yaml and helm-chart.yaml:
git config user.name "docker-bot"
git config user.email "docker-bot@users.noreply.github.com"

# Result: sync-to-dev.yaml won't trigger on bot commits
# But PR is still created for human review and approval
```

**Why We Keep PRs in docker.yaml and helm-chart.yaml**:
1. **Human oversight**: Changes reviewed before production
2. **Audit trail**: Clear history of what changed and why
3. **Safety net**: Prevents automatic deployment of broken configs
4. **Compliance**: Meets review requirements for production deployments

**Architecture Benefits**:
- 🎯 **Separation of concerns**: Three distinct workflows for different triggers
- 🔒 **Loop prevention**: Bot account prevents cyclic executions
- 👥 **Human-in-loop**: PRs ensure review before environment promotion
- 🚀 **Fast feedback**: Immediate builds and updates to feature branch
- 🌐 **Multi-environment**: Simultaneous updates to dev, uat, prod overlays
- 📦 **Artifact traceability**: Every build tagged with commit SHA and chart version

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
5. **Intelligent workflow orchestration** with loop prevention
6. **Pull request automation** with human review gates

This project serves as a **comprehensive template** for building production-grade CI/CD pipelines for modern web applications, demonstrating enterprise-level practices in a practical, hands-on implementation.

## 📚 Explore the Documentation

The complete documentation is organized into logical sections:

1. **[Getting Started](docs/getting-started.md)** - Installation, setup, and first deployment
2. **[Project Structure](docs/project-structure.md)** - Understanding the codebase organization
3. **[Deployment Scripts](docs/deployment-scripts.md)** - Deep dive into automation scripts
4. **[CI/CD Workflows](docs/cicd-workflows.md)** - GitHub Actions pipeline details
5. **[Configuration](docs/configuration.md)** - Environment and resource management
6. **[Monitoring](docs/monitoring.md)** - Debugging and troubleshooting guide
7. **[Production Features](docs/production-features.md)** - Production-ready capabilities
8. **[Learning Outcomes](docs/learning-outcomes.md)** - Skills and knowledge gained

### View Documentation

```bash
# Install MkDocs with Material theme
pip install mkdocs-material

# Serve documentation locally
mkdocs serve

# Build static documentation site
mkdocs build
```

Access the documentation at `http://127.0.0.1:8000` with full search, navigation, and syntax highlighting.

## 🤝 Contributing

Contributions are welcome! Please refer to the documentation for architecture details and best practices.

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.



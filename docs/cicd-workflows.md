# CI/CD Workflows

Comprehensive guide to the GitHub Actions workflows that power the automated CI/CD pipeline.

## Overview

This project implements a **sophisticated three-workflow GitOps pipeline** that handles application builds, chart packaging, and configuration synchronization with intelligent pull request management.

## Workflow Architecture

```
Code/Config Change → GitHub Push to Feature Branch
          ↓
┌─────────────────────┬─────────────────────┬─────────────────────┐
│   docker.yaml       │   helm-chart.yaml   │   sync-to-dev.yaml  │
│   Build & Deploy    │   Package & Publish │   Config Sync       │
└─────────────────────┴─────────────────────┴─────────────────────┘
          ↓                     ↓                     ↓
    Pull Request to dev → Manual Review → Merge → ArgoCD Deploys
```

## Workflow 1: Docker Image Pipeline

**File**: `.github/workflows/docker.yaml`

### Purpose
Automatically builds and publishes container images whenever application code changes, then updates environment configurations and creates pull requests.

### Trigger Conditions

```yaml
on:
  push:
    branches: ["feature"]
    paths:
      - 'src/**'
      - 'Dockerfile'
      - 'go.mod'
      - 'go.sum'
      - '.github/workflows/docker.yaml'
  workflow_dispatch: {}
```

**Triggers on:**
- Changes to application source code (`src/**`)
- Dockerfile modifications
- Go dependency changes (`go.mod`, `go.sum`)
- Workflow file changes
- Manual workflow dispatch

### Workflow Steps

#### Phase 1: Build and Push

```yaml
steps:
  - name: Checkout repository
    uses: actions/checkout@v4
  
  - name: Login to Docker Hub
    uses: docker/login-action@v3
    with:
      username: ${{ vars.DOCKERHUB_USERNAME }}
      password: ${{ secrets.DOCKERHUB_TOKEN }}
  
  - name: Set up QEMU
    uses: docker/setup-qemu-action@v3
  
  - name: Set up Docker Buildx
    uses: docker/setup-buildx-action@v3
  
  - name: Build and push Docker image
    uses: docker/build-push-action@v6
    with:
      context: .
      file: Dockerfile
      platforms: linux/amd64,linux/arm64
      push: true
      tags: |
        abhishek1009/todo-app:latest
        abhishek1009/todo-app:${{ github.sha }}
      cache-from: type=registry,ref=abhishek1009/todo-app:cache
      cache-to: type=registry,ref=abhishek1009/todo-app:cache,mode=max
```

**What happens:**
- ✅ Multi-architecture build (amd64/arm64)
- ✅ Docker Hub authentication
- ✅ Build caching for performance
- ✅ Dual tagging: `latest` + commit SHA

#### Phase 2: Update ArgoCD Configurations

```yaml
- name: Update image tags in overlays
  run: |
    # Update dev environment
    sed -i "s|image:.*|image: abhishek1009/todo-app:${{ github.sha }}|g" \
      argocd/apps/todo-app/overlays/dev/values.yaml
    
    # Update uat environment
    sed -i "s|image:.*|image: abhishek1009/todo-app:${{ github.sha }}|g" \
      argocd/apps/todo-app/overlays/uat/values.yaml
    
    # Update prod environment
    sed -i "s|image:.*|image: abhishek1009/todo-app:${{ github.sha }}|g" \
      argocd/apps/todo-app/overlays/prod/values.yaml
```

**What happens:**
- Updates image tags in all environment overlays
- Uses commit SHA for traceability
- Maintains separate configurations per environment

#### Phase 3: Commit and Push

```yaml
- name: Commit and push changes
  run: |
    git config user.name "docker-bot"
    git config user.email "docker-bot@users.noreply.github.com"
    git add argocd/apps/todo-app/overlays/*/values.yaml
    git commit -m "Update image tag to ${{ github.sha }}"
    git push origin feature
```

**Why Docker Bot?**
- Prevents infinite workflow loops
- GitHub won't trigger workflows on bot commits
- Maintains clean commit history

#### Phase 4: Create Pull Request

```yaml
- name: Create Pull Request
  uses: peter-evans/create-pull-request@v5
  with:
    token: ${{ secrets.GITHUB_TOKEN }}
    branch: feature
    base: dev
    title: "Update Docker image to ${{ github.sha }}"
    labels: |
      automated-pr
      docker-update
```

**What happens:**
- Creates/updates PR from feature to dev
- Adds descriptive labels
- Includes commit SHA in title

## Workflow 2: Helm Chart Pipeline

**File**: `.github/workflows/helm-chart.yaml`

### Purpose
Packages and publishes Helm charts to OCI registry when chart configurations change, then updates Kustomize references.

### Trigger Conditions

```yaml
on:
  push:
    branches: ["feature"]
    paths:
      - 'helm/**'
      - '.github/workflows/helm-chart.yaml'
  workflow_dispatch: {}
```

**Triggers on:**
- Changes to Helm chart files (`helm/**`)
- Workflow file modifications
- Manual workflow dispatch

### Workflow Steps

#### Phase 1: Package and Publish

```yaml
- name: Build and push helm chart (stable version)
  run: |
    # Extract chart metadata
    CHART_NAME=$(grep 'name:' helm/Chart.yaml | awk '{print $2}')
    CHART_NAME=${CHART_NAME//\"/}
    
    CHART_VERSION=$(grep 'version:' helm/Chart.yaml | awk '{print $2}')
    CHART_VERSION=${CHART_VERSION//\"/}
    
    # Package and push
    helm package helm --destination ./packages
    helm push ./packages/$CHART_NAME-$CHART_VERSION.tgz \
      oci://registry-1.docker.io/abhishek1009
```

**What happens:**
- Extracts chart name and version from `Chart.yaml`
- Packages Helm chart
- Pushes to Docker Hub OCI registry

#### Phase 2: Development Version

```yaml
- name: Build and push helm chart (commit SHA version)
  run: |
    # Create version with commit SHA
    CHART_VERSION_FROM_COMMIT=${GITHUB_SHA::7}
    CHART_VERSION_FROM_COMMIT="$CHART_VERSION-$CHART_VERSION_FROM_COMMIT"
    
    # Package with custom version
    helm package helm --version $CHART_VERSION_FROM_COMMIT \
      --destination ./packages
    
    # Push development version
    helm push ./packages/$CHART_NAME-$CHART_VERSION_FROM_COMMIT.tgz \
      oci://registry-1.docker.io/abhishek1009
```

**Why two versions?**
- **Stable version** (1.0.0): For releases
- **Development version** (1.0.0-a1b2c3d): For tracking changes

#### Phase 3: Update Kustomization Files

```yaml
- name: Update chart versions in overlays
  run: |
    # Update dev kustomization
    sed -i "s|version:.*|version: $CHART_VERSION|g" \
      argocd/apps/todo-app/overlays/dev/kustomization.yaml
    
    # Update uat kustomization
    sed -i "s|version:.*|version: $CHART_VERSION|g" \
      argocd/apps/todo-app/overlays/uat/kustomization.yaml
    
    # Update prod kustomization
    sed -i "s|version:.*|version: $CHART_VERSION|g" \
      argocd/apps/todo-app/overlays/prod/kustomization.yaml
```

**What happens:**
- Updates Helm chart version references
- Synchronizes all environments
- Ensures consistent versioning

#### Phase 4: Commit and Create PR

Similar to docker.yaml, commits changes using bot account and creates PR.

## Workflow 3: ArgoCD Sync Pipeline

**File**: `.github/workflows/sync-to-dev.yaml`

### Purpose
Detects any changes in ArgoCD configurations and creates PRs to synchronize feature branch to dev.

### Trigger Conditions

```yaml
on:
  push:
    branches: ["feature"]
    paths:
      - 'argocd/**'
  workflow_dispatch: {}
```

**Triggers on:**
- Any changes to `argocd/` directory
- Manual workflow dispatch

### Workflow Steps

```yaml
- name: Create Pull Request
  uses: peter-evans/create-pull-request@v5
  with:
    token: ${{ secrets.GITHUB_TOKEN }}
    branch: feature
    base: dev
    title: "Sync ArgoCD configs from feature to dev"
    body: |
      Automated PR to sync ArgoCD configuration changes.
      
      Changed files:
      ${{ steps.changes.outputs.files }}
    labels: |
      automated-pr
      argocd-sync
```

**What happens:**
- Creates PR for any ArgoCD changes
- Lists changed files in PR description
- Adds descriptive labels

## Preventing Infinite Loops

### The Problem
```
docker.yaml runs → updates argocd/ → commits & pushes
    ↓
sync-to-dev.yaml detects change → creates PR
    ↓
PR merged → triggers docker.yaml again → LOOP! ❌
```

### The Solution

#### Bot Account Commits
```yaml
git config user.name "docker-bot"
git config user.email "docker-bot@users.noreply.github.com"
```

#### GitHub Actions Behavior
- GitHub won't trigger workflows on bot commits
- Prevents cyclic executions
- Maintains automation safety

### Why Keep PRs?

Even though sync-to-dev.yaml won't trigger on bot commits, we still create PRs because:

1. **Human oversight**: Review before production
2. **Audit trail**: Clear change history
3. **Safety net**: Catch configuration errors
4. **Compliance**: Meet review requirements

## Complete CI/CD Flow

```
Developer pushes code to feature branch
          ↓
┌─────────────────────────────────────────────────────────┐
│ GitHub Actions detect changes and trigger workflows    │
└─────────────────────────────────────────────────────────┘
          ↓
┌─────────────────────┬─────────────────────┬─────────────────────┐
│ docker.yaml         │ helm-chart.yaml     │ sync-to-dev.yaml    │
│ • Build image       │ • Package chart     │ • Detect changes    │
│ • Push to registry  │ • Push to OCI       │ • Create PR         │
│ • Update values     │ • Update kustomize  │                     │
│ • Commit (bot)      │ • Commit (bot)      │                     │
│ • Push to feature   │ • Push to feature   │                     │
│ • Create PR         │ • Create PR         │                     │
└─────────────────────┴─────────────────────┴─────────────────────┘
          ↓
    Pull Request Created/Updated (feature → dev)
          ↓
    Manual Review & Approval Required
          ↓
    Merge to Dev Branch
          ↓
    ArgoCD Detects Change in Dev Branch
          ↓
    Pulls Latest Artifacts (images + charts)
          ↓
    Deploys to Kubernetes Cluster (dev environment)
          ↓
    Repeat process for UAT → Prod promotion
```

## Environment Variables & Secrets

### Required GitHub Variables
```yaml
DOCKERHUB_USERNAME: Your Docker Hub username
APP_IMAGE_NAME: Name of your application image
HELM_REGISTRY: OCI registry URL
```

### Required GitHub Secrets
```yaml
DOCKERHUB_TOKEN: Docker Hub access token
GITHUB_TOKEN: Automatically provided by GitHub Actions
```

## Best Practices

### Tagging Strategy
- **`latest`**: Most recent build
- **`${GITHUB_SHA}`**: Specific commit traceability
- **`${CHART_VERSION}-${SHA}`**: Development chart versions

### Build Caching
```yaml
cache-from: type=registry,ref=abhishek1009/todo-app:cache
cache-to: type=registry,ref=abhishek1009/todo-app:cache,mode=max
```
- Speeds up subsequent builds
- Reduces build time by 50-70%
- Shares cache across workflows

### Multi-Architecture Support
```yaml
platforms: linux/amd64,linux/arm64
```
- Supports Intel and ARM processors
- Works on diverse Kubernetes clusters
- Future-proof for Apple Silicon, AWS Graviton

### Workflow Permissions
```yaml
permissions:
  contents: write
  pull-requests: write
```
- Minimal required permissions
- Follows principle of least privilege
- Ensures security compliance

## Troubleshooting

### Workflow Not Triggering
- Check path filters match your changes
- Verify branch name is correct
- Check if previous workflow is running

### Build Failures
- Check Docker Hub credentials
- Verify Go dependencies are valid
- Review build logs for specific errors

### PR Not Created
- Verify GITHUB_TOKEN has permissions
- Check if PR already exists
- Review workflow logs

## Next Steps

- [Configure environments](configuration.md)
- [Monitor deployments](monitoring.md)
- [Learn about production features](production-features.md)

# Deployment Scripts

Comprehensive guide to the automated deployment scripts in the root directory.

## Overview

Three powerful scripts automate the entire ArgoCD deployment and configuration process:

- **`application.yaml`** - Root application template
- **`argo-deploy.sh`** - Main deployment automation
- **`enable-helm.sh`** - Helm integration setup

## `application.yaml` - ArgoCD Root Application

### Purpose
Implements the **"App of Apps" pattern** in ArgoCD, allowing management of multiple applications across different environments from a single root application.

### Template Structure
```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: root-application
  namespace: argocd
spec:
  project: default
  destination:
    server: https://kubernetes.default.svc
    namespace: argocd
  source:
    repoURL: https://github.com/abhishekghoshh/todo-app-cicd
    targetRevision: ${targetEnvironment}
    path: argocd/argo-apps/${targetEnvironment}
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
```

### Key Features

#### Environment Variable Substitution
```bash
# Uses ${targetEnvironment} for dynamic targeting
export targetEnvironment=feature
envsubst '${targetEnvironment}' < application.yaml | kubectl apply -f -
```

#### Dynamic Branching
- Points to different Git branches based on environment
- Each environment has its own branch (feature, dev, uat, prod)
- ArgoCD tracks branch-specific changes

#### Automated Sync Policy
- **prune**: Removes resources deleted from Git
- **selfHeal**: Automatically corrects drift

## `argo-deploy.sh` - Automated Deployment Script

### Synopsis
```bash
./argo-deploy.sh <environment>
./argo-deploy.sh destroy  # Complete cleanup
```

### Usage Examples
```bash
# Deploy to feature environment
./argo-deploy.sh feature

# Deploy to development
./argo-deploy.sh dev

# Deploy to production
./argo-deploy.sh prod

# Remove everything
./argo-deploy.sh destroy
```

### Complete Workflow

#### Step 1: Validation
```bash
if [ -z "$1" ]; then
  echo "❌ Please provide an environment: dev | stage | prod"
  exit 1
fi
```
- Validates environment argument
- Provides usage instructions
- Handles destroy command

#### Step 2: ArgoCD Installation Check
```bash
if kubectl get namespace argocd >/dev/null 2>&1 && \
   kubectl get deployment argocd-server -n argocd >/dev/null 2>&1; then
    echo "✅ ArgoCD is already installed"
else
    # Install ArgoCD
fi
```
- Checks if ArgoCD namespace exists
- Verifies argocd-server deployment
- Skips installation if already present

#### Step 3: Fresh Installation
```bash
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
kubectl wait deployment/argocd-server -n argocd \
  --for=condition=Available=True --timeout=180s
```
- Creates argocd namespace
- Applies official ArgoCD manifests
- Waits for server readiness

#### Step 4: Helm Integration
```bash
./enable-helm.sh
```
- Calls enable-helm.sh script
- Configures Helm support in Kustomize
- Required for Helm chart deployment

#### Step 5: Root Application Deployment
```bash
export targetEnvironment="$TARGET_ENV"
envsubst '${targetEnvironment}' < "$ROOT_FILE" > "$TEMP_FILE"
kubectl apply -f "$TEMP_FILE"
rm "$TEMP_FILE"
```
- Substitutes environment variable
- Generates temporary rendered file
- Applies to Kubernetes
- Cleans up temporary file

#### Step 6: Port-Forwarding & Access
```bash
# Display admin password
kubectl -n argocd get secret argocd-initial-admin-secret \
  -o jsonpath="{.data.password}" | base64 -d

# Start port-forwarding
kubectl port-forward svc/argocd-server -n argocd 8080:443
```
- Retrieves admin password
- Starts automatic port-forwarding
- Provides access instructions

### Destroy Mode

Complete cleanup of ArgoCD and all applications:

```bash
./argo-deploy.sh destroy
```

**What it does:**
1. Deletes all ArgoCD applications
2. Deletes all app projects
3. Waits for graceful cleanup
4. Removes ArgoCD installation
5. Deletes argocd namespace

## `enable-helm.sh` - Helm Integration Setup

### Purpose
Configures ArgoCD to support Helm charts within Kustomize builds, enabling the hybrid Kustomize+Helm approach.

### Workflow

#### Step 1: Ensure Namespace
```bash
kubectl create namespace argocd 2>/dev/null || true
```
- Creates namespace if missing
- Ignores error if already exists

#### Step 2: Patch ConfigMap
```bash
kubectl patch configmap argocd-cm -n argocd \
  --type merge \
  -p '{"data":{"kustomize.buildOptions":"--enable-helm --load-restrictor LoadRestrictionsNone"}}'
```
- Patches existing argocd-cm ConfigMap
- Adds Helm build options
- Enables remote chart loading

#### Step 3: Create ConfigMap if Missing
```bash
kubectl apply -n argocd -f - <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-cm
  namespace: argocd
data:
  kustomize.buildOptions: "--enable-helm --load-restrictor LoadRestrictionsNone"
EOF
```
- Fallback if ConfigMap doesn't exist
- Creates new ConfigMap with Helm support

#### Step 4: Restart Components
```bash
kubectl rollout restart deployment argocd-repo-server -n argocd
kubectl rollout restart deployment argocd-server -n argocd
kubectl rollout restart deployment argocd-applicationset-controller -n argocd
kubectl rollout restart deployment argocd-dex-server -n argocd
```
- Restarts all ArgoCD components
- Applies new configuration
- Ensures changes take effect

#### Step 5: Wait for Readiness
```bash
kubectl wait deployment/argocd-repo-server \
  -n argocd --for=condition=Available=True --timeout=120s
kubectl wait deployment/argocd-server \
  -n argocd --for=condition=Available=True --timeout=120s
```
- Waits for components to be ready
- Ensures stability before proceeding

### Build Options Explained

#### `--enable-helm`
- Enables Helm chart support in Kustomize
- Allows helmCharts in kustomization.yaml
- Processes Helm templates

#### `--load-restrictor LoadRestrictionsNone`
- Removes file loading restrictions
- Allows loading from any directory
- Required for remote Helm charts

## Best Practices

### Error Handling
All scripts include:
- Exit on error (`set -e`)
- Error messages with emojis
- Graceful fallbacks
- Status confirmations

### Idempotency
Scripts can be run multiple times safely:
- Check before install
- Skip if already configured
- Update instead of create

### User Feedback
Clear console output:
- ✅ Success indicators
- ❌ Error messages
- ⏳ Progress indicators
- 🎉 Completion messages

### Security
- Displays admin password securely
- Uses kubectl port-forward
- No hardcoded credentials

## Common Use Cases

### Initial Setup
```bash
# One-time setup for new cluster
./argo-deploy.sh feature
```

### Environment Switch
```bash
# Switch from feature to dev
./argo-deploy.sh dev
```

### Troubleshooting
```bash
# Reset ArgoCD completely
./argo-deploy.sh destroy
./argo-deploy.sh feature
```

### Manual Helm Enable
```bash
# If Helm support needs reconfiguration
./enable-helm.sh
```

## Next Steps

- [Configure CI/CD workflows](cicd-workflows.md)
- [Manage environments](configuration.md)
- [Monitor deployments](monitoring.md)

# Configuration Guide

Learn how to configure and customize the Todo App CI/CD pipeline for different environments.

## Environment Management

The project supports multiple environments with isolated configurations:

- **feature**: Development and testing
- **dev**: Development environment
- **uat**: User acceptance testing
- **prod**: Production environment

## Configuration Structure

### Directory Layout
```
argocd/apps/todo-app/
├── base/
│   ├── kustomization.yaml      # Base Kustomize config
│   └── values.yaml             # Base Helm values
└── overlays/
    ├── dev/
    │   ├── kustomization.yaml  # Dev-specific Kustomize
    │   └── values.yaml         # Dev-specific Helm values
    ├── feature/
    │   ├── kustomization.yaml
    │   └── values.yaml
    ├── prod/
    │   ├── kustomization.yaml
    │   └── values.yaml
    └── uat/
        ├── kustomization.yaml
        └── values.yaml
```

## Base Configuration

### `base/values.yaml`

Common configuration shared across all environments:

```yaml
nameOverride: ""

deployment:
  app:
    replicas: 1
    name: "todo-app"
    image:
      host: ""
      name: "abhishek1009/todo-app"
      tag: "latest"
      pullPolicy: IfNotPresent
    port: 8080
    
    env:
      sessionTtlMinutes: "10"
    
    expose:
      enabled: true
      ingressType: "nginx"
      host: "app.todo.local"
      annotations:
        nginx.ingress.kubernetes.io/proxy-body-size: "20m"
    
    pvc:
      size: 1Gi
      storageClassName: "standard"
    
    resources:
      requests:
        cpu: "200m"
        memory: "256Mi"
      limits:
        cpu: "1000m"
        memory: "1024Mi"
  
  db:
    replicas: 1
    name: "mongodb"
    image:
      name: "mongo"
      tag: "latest"
    port: 27017
    
    expose:
      enabled: false
    
    pvc:
      size: 8Gi
      storageClassName: "standard"
    
    resources:
      requests:
        cpu: "200m"
        memory: "256Mi"
      limits:
        cpu: "1000m"
        memory: "1024Mi"
```

### `base/kustomization.yaml`

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: todo-app

resources: []
```

## Environment Overlays

### Feature Environment

**`overlays/feature/values.yaml`**:
```yaml
nameOverride: "todo-app"

deployment:
  app:
    expose:
      host: "app.todo.local"
  db:
    expose:
      host: "mongo.local"
```

**`overlays/feature/kustomization.yaml`**:
```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: todo-app

resources:
- ../../base

helmCharts:
  - name: todo-app-chart
    repo: oci://registry-1.docker.io/abhishek1009
    version: 1.0.0
    namespace: todo-app
    valuesFile: ../../base/values.yaml
    additionalValuesFiles:
      - values.yaml
    valuesInline:
      nameOverride: todo-app
```

### Development Environment

**`overlays/dev/values.yaml`**:
```yaml
nameOverride: "todo-app-dev"

deployment:
  app:
    image:
      tag: "latest"  # Updated by CI/CD
    expose:
      host: "app-todo-dev.local"
  db:
    expose:
      host: "mongo-dev.local"
```

### Production Environment

**`overlays/prod/values.yaml`**:
```yaml
nameOverride: "todo-app-prod"

deployment:
  app:
    replicas: 3  # Higher availability
    image:
      tag: "1.0.0"  # Stable version
    expose:
      host: "app.todo.com"
      annotations:
        cert-manager.io/cluster-issuer: "letsencrypt-prod"
    resources:
      requests:
        cpu: "500m"
        memory: "512Mi"
      limits:
        cpu: "2000m"
        memory: "2048Mi"
  
  db:
    pvc:
      size: 100Gi  # Larger storage
    resources:
      requests:
        cpu: "1000m"
        memory: "2048Mi"
      limits:
        cpu: "4000m"
        memory: "4096Mi"
```

## Configuration Options

### Application Settings

#### Image Configuration
```yaml
image:
  host: ""                          # Optional registry host
  name: "abhishek1009/todo-app"     # Image name
  tag: "latest"                     # Image tag (updated by CI/CD)
  pullPolicy: IfNotPresent          # Pull policy
```

#### Replica Count
```yaml
replicas: 3  # Number of pod replicas
```

#### Environment Variables
```yaml
env:
  sessionTtlMinutes: "30"  # Session timeout in minutes
```

#### Port Configuration
```yaml
port: 8080  # Application port
```

### Ingress Configuration

#### Basic Ingress
```yaml
expose:
  enabled: true
  ingressType: "nginx"
  host: "app.todo.com"
```

#### Advanced Ingress with TLS
```yaml
expose:
  enabled: true
  ingressType: "nginx"
  host: "app.todo.com"
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/proxy-body-size: "50m"
```

### Storage Configuration

#### Persistent Volume Claims
```yaml
pvc:
  size: 10Gi                    # Storage size
  storageClassName: "standard"  # Storage class
```

**Common storage classes:**
- `standard`: Default storage
- `gp2`, `gp3`: AWS EBS
- `premium-rwo`: GCP persistent disk
- `azure-disk`: Azure disk storage

### Resource Management

#### Development Resources
```yaml
resources:
  requests:
    cpu: "100m"
    memory: "128Mi"
  limits:
    cpu: "500m"
    memory: "512Mi"
```

#### Production Resources
```yaml
resources:
  requests:
    cpu: "1000m"
    memory: "1024Mi"
  limits:
    cpu: "4000m"
    memory: "4096Mi"
```

### MongoDB Configuration

#### Basic MongoDB
```yaml
db:
  replicas: 1
  name: "mongodb"
  image:
    name: "mongo"
    tag: "7.0"
  port: 27017
  pvc:
    size: 20Gi
```

#### Production MongoDB
```yaml
db:
  replicas: 3  # ReplicaSet for HA
  name: "mongodb-prod"
  image:
    name: "mongo"
    tag: "7.0"
  port: 27017
  pvc:
    size: 100Gi
    storageClassName: "fast-ssd"
  resources:
    requests:
      cpu: "2000m"
      memory: "4096Mi"
    limits:
      cpu: "4000m"
      memory: "8192Mi"
```

## Helm Chart Configuration

### Chart Metadata (`helm/Chart.yaml`)

```yaml
apiVersion: v2
name: todo-app-chart
description: A Helm chart for the Go-HTMX-MongoDB Todo App
version: "1.0.0"      # Updated manually
appVersion: "1.0.0"   # Application version
```

### Version Management

The `version` field in `Chart.yaml` controls:
- Helm chart versioning
- OCI registry tags
- ArgoCD chart references

**Updating version:**
1. Edit `helm/Chart.yaml`
2. Increment version number
3. Commit and push
4. CI/CD automatically packages and publishes

## ArgoCD Application Configuration

### Application Definition (`argocd/argo-apps/feature/todo-app.yaml`)

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: todo-app
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://github.com/abhishekghoshh/todo-app-cicd
    targetRevision: feature
    path: argocd/apps/todo-app/overlays/feature
    kustomize:
      enableHelm: true
  destination:
    server: https://kubernetes.default.svc
    namespace: todo-app
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
    - CreateNamespace=true
```

### Sync Policy Options

#### Automated Sync
```yaml
syncPolicy:
  automated:
    prune: true      # Remove deleted resources
    selfHeal: true   # Auto-correct drift
```

#### Manual Sync
```yaml
syncPolicy:
  automated: {}  # Empty for manual sync
```

#### Sync Options
```yaml
syncOptions:
  - CreateNamespace=true    # Auto-create namespace
  - Validate=true           # Validate manifests
  - PruneLast=true         # Prune after other resources
```

## Environment Variables

### Application Environment Variables

Set in deployment manifests:

```yaml
env:
  - name: MONGO_URI
    value: "mongodb://mongodb-svc:27017"
  - name: SESSION_TTL_MINUTES
    value: "30"
  - name: LOG_LEVEL
    value: "info"
```

### Using Secrets

```yaml
env:
  - name: MONGO_PASSWORD
    valueFrom:
      secretKeyRef:
        name: mongo-secret
        key: password
```

## Best Practices

### 1. Environment Separation
- Keep production config separate from dev
- Use different resource limits
- Configure appropriate storage sizes

### 2. Version Control
- Tag stable releases
- Use semantic versioning
- Document breaking changes

### 3. Resource Planning
- Start with conservative limits
- Monitor actual usage
- Adjust based on metrics

### 4. Security
- Never commit secrets
- Use Kubernetes secrets
- Enable RBAC and network policies

### 5. Storage Management
- Choose appropriate storage classes
- Plan for data growth
- Configure backup solutions

## Customization Examples

### Custom Domain with TLS

```yaml
deployment:
  app:
    expose:
      enabled: true
      host: "todos.example.com"
      annotations:
        cert-manager.io/cluster-issuer: "letsencrypt-prod"
        nginx.ingress.kubernetes.io/force-ssl-redirect: "true"
```

### Custom Resource Limits

```yaml
deployment:
  app:
    resources:
      requests:
        cpu: "250m"
        memory: "256Mi"
      limits:
        cpu: "1000m"
        memory: "1024Mi"
```

### Multiple Replicas

```yaml
deployment:
  app:
    replicas: 5
```

### External MongoDB

```yaml
deployment:
  app:
    env:
      mongoUri: "mongodb://external-mongo.example.com:27017"
  
  db:
    enabled: false  # Disable internal MongoDB
```

## Troubleshooting Configuration

### Configuration Not Applied
```bash
# Check ArgoCD sync status
argocd app get todo-app

# Force sync
argocd app sync todo-app --force
```

### Invalid YAML
```bash
# Validate kustomization
kustomize build argocd/apps/todo-app/overlays/dev

# Validate Helm values
helm template helm -f argocd/apps/todo-app/overlays/dev/values.yaml
```

### Resource Issues
```bash
# Check pod resources
kubectl top pods -n todo-app

# Describe pod for events
kubectl describe pod <pod-name> -n todo-app
```

## Next Steps

- [Monitor your deployment](monitoring.md)
- [Learn about production features](production-features.md)
- [Understand CI/CD workflows](cicd-workflows.md)

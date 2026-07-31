# Project Structure

Understanding the organization of the Todo App CI/CD pipeline.

## Directory Layout

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
│       ├── mongo/              # MongoDB deployment resources
│       └── todo/               # Todo app deployment resources
├── argocd/                     # GitOps configuration
│   ├── apps/todo-app/          # Application-specific configs
│   │   ├── base/               # Base Kustomize configuration
│   │   └── overlays/           # Environment-specific overrides
│   │       ├── dev/            # Development environment
│   │       ├── feature/        # Feature environment
│   │       ├── prod/           # Production environment
│   │       └── uat/            # UAT environment
│   └── argo-apps/              # ArgoCD Application definitions
│       ├── dev/                # Dev app definitions
│       ├── feature/            # Feature app definitions
│       ├── prod/               # Prod app definitions
│       └── uat/                # UAT app definitions
├── .github/workflows/          # CI/CD automation
│   ├── docker.yaml             # Docker image build and push
│   ├── helm-chart.yaml         # Helm chart packaging and publish
│   └── sync-to-dev.yaml        # ArgoCD config sync automation
├── docs/                       # Documentation (you are here!)
├── application.yaml            # ArgoCD root application
├── argo-deploy.sh              # Automated deployment script
├── enable-helm.sh              # Helm integration setup
├── docker-compose.yaml         # Local development setup
├── Dockerfile                  # Multi-stage container build
└── version                     # Version tracking file
```

## Core Components

### Application Code (`src/`)

The Go application is organized into clean, modular components:

- **`main.go`**: Entry point, HTTP server setup, route definitions
- **`handlers.go`**: HTTP request handlers, middleware, authentication
- **`models.go`**: Data structures (User, Todo, Session)
- **`db.go`**: MongoDB connection, CRUD operations, session management

### Helm Charts (`helm/`)

Production-ready Kubernetes manifests packaged as Helm charts:

- **`Chart.yaml`**: Chart metadata and version information
- **`values.yaml`**: Default configuration values
- **`templates/`**: Kubernetes resource templates
  - Deployments for app and MongoDB
  - Services for internal communication
  - Ingress for external access
  - PersistentVolumeClaims for data storage

### GitOps Configuration (`argocd/`)

Two-level structure for ArgoCD management:

#### Application Configs (`apps/todo-app/`)
- **`base/`**: Common configuration shared across environments
- **`overlays/`**: Environment-specific overrides using Kustomize
  - Each overlay contains `values.yaml` and `kustomization.yaml`
  - Customizes image tags, resource limits, ingress hosts, etc.

#### ArgoCD Apps (`argo-apps/`)
- Application definitions for each environment
- Points to specific Git branch and overlay path
- Defines sync policies (automated, manual, self-heal)

### CI/CD Workflows (`.github/workflows/`)

Three intelligent workflows that work together:

1. **`docker.yaml`**: Builds container images on code changes
2. **`helm-chart.yaml`**: Packages Helm charts on configuration changes
3. **`sync-to-dev.yaml`**: Creates PRs for environment promotion

## Root-Level Scripts

### `application.yaml`
ArgoCD root application implementing the "App of Apps" pattern:
- Manages all environment applications from one place
- Uses environment variable substitution
- Enables dynamic environment targeting

### `argo-deploy.sh`
One-command deployment automation:
- Installs/configures ArgoCD
- Deploys applications to any environment
- Handles cleanup and destruction

### `enable-helm.sh`
Helm integration for ArgoCD:
- Patches ArgoCD configuration
- Enables Helm support in Kustomize
- Restarts necessary components

## Environment Flow

```
Development Workflow:
feature branch → docker.yaml → builds image → updates overlays → PR to dev
                ↓
             helm-chart.yaml → packages chart → updates kustomization → PR to dev
                ↓
             sync-to-dev.yaml → detects changes → creates PR
                ↓
             Manual merge to dev → ArgoCD syncs → Deploys to dev environment
                ↓
             Promote to UAT → Promote to Prod
```

## Key Design Decisions

### Kustomize + Helm Hybrid
- **Base**: Common Helm values
- **Overlays**: Environment-specific customizations
- **Benefits**: DRY principle, flexible configuration

### OCI Registry for Both Images and Charts
- **Docker Hub**: Stores both container images and Helm charts
- **Benefits**: Single registry, simplified authentication

### GitOps with ArgoCD
- **Declarative**: Git as single source of truth
- **Automated**: Self-healing and auto-sync
- **Auditable**: Complete change history

### Three-Workflow CI/CD
- **Separation of concerns**: Different triggers for different artifacts
- **Loop prevention**: Bot account prevents infinite triggers
- **Human oversight**: PRs require approval before promotion

## Next Steps

- [Understand deployment scripts](deployment-scripts.md)
- [Configure CI/CD workflows](cicd-workflows.md)
- [Manage environments](configuration.md)

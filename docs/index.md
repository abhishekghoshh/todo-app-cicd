# Todo App CI/CD Pipeline

A comprehensive **GitOps-powered CI/CD pipeline** for a modern Go-based Todo application with complete Kubernetes deployment automation using ArgoCD, Helm, and GitHub Actions.

## 🚀 Overview

This project demonstrates a **production-ready CI/CD pipeline** that automatically builds, tests, packages, and deploys a Go web application to Kubernetes using modern cloud-native practices.

## Technology Stack

- **Backend**: Go 1.25+ with native HTTP routing, MongoDB integration
- **Frontend**: HTMX-powered interactive UI with server-side rendering
- **Database**: MongoDB with session management and TTL indexes
- **Containerization**: Multi-stage Docker builds with Alpine Linux
- **Orchestration**: Kubernetes with Helm charts and Kustomize overlays
- **GitOps**: ArgoCD for automated deployment and management
- **CI/CD**: GitHub Actions for build automation and artifact publishing

## Key Features

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

## Architecture

### Application Architecture
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTMX Client   │────│   Go Backend    │────│    MongoDB      │
│  (Dashboard)    │    │  (REST + Auth)  │    │  (Persistence)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### CI/CD Pipeline Flow
```
┌─────────────┐   ┌──────────────┐   ┌─────────────┐   ┌─────────────┐
│ Code Push   │──▶│ GitHub       │──▶│ ArgoCD      │──▶│ Kubernetes  │
│ (feature)   │   │ Actions      │   │ Sync        │   │ Deployment  │
│             │   │ (Build/Push) │   │ (GitOps)    │   │ (Live App)  │
└─────────────┘   └──────────────┘   └─────────────┘   └─────────────┘
```

## Quick Links

- [Getting Started](getting-started.md) - Installation and setup
- [Project Structure](project-structure.md) - Directory layout and organization
- [Deployment Scripts](deployment-scripts.md) - Automated deployment tools
- [CI/CD Workflows](cicd-workflows.md) - GitHub Actions automation
- [Configuration](configuration.md) - Environment management
- [Monitoring](monitoring.md) - Troubleshooting and debugging
- [Production Features](production-features.md) - Production-ready capabilities
- [Learning Outcomes](learning-outcomes.md) - What this project teaches

## Prerequisites

- Docker installed and running
- Kubernetes cluster (local or cloud)
- kubectl configured and connected to cluster
- Helm 3.x installed
- Git repository access
- Docker Hub account (for image and chart registry)

## What You'll Achieve

By following this documentation, you'll set up a complete CI/CD pipeline that:

1. **Automatically builds** Docker images on code changes
2. **Packages Helm charts** and publishes to OCI registry
3. **Manages multiple environments** (dev, staging, production)
4. **Deploys via GitOps** using ArgoCD
5. **Provides automated rollbacks** and self-healing
6. **Implements security best practices** for production workloads

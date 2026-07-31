# Production Features

Comprehensive overview of production-ready features and best practices implemented in this project.

## High Availability

### Application Redundancy

**Multiple Replicas**:
```yaml
deployment:
  app:
    replicas: 3  # Run multiple instances
```

**Benefits**:
- Zero-downtime deployments
- Load distribution
- Fault tolerance
- Better resource utilization

**Rolling Updates**:
```bash
# Kubernetes automatically performs rolling updates
kubectl set image deployment/todo-app \
  app=abhishek1009/todo-app:new-version -n todo-app
```

### Database High Availability

**MongoDB ReplicaSet** (Production):
```yaml
deployment:
  db:
    replicas: 3  # Primary + 2 secondaries
    pvc:
      size: 100Gi
    resources:
      requests:
        cpu: "2000m"
        memory: "4096Mi"
```

**Features**:
- Automatic failover
- Data redundancy
- Read scaling
- Point-in-time recovery

### Load Balancing

**Kubernetes Service**:
- Automatic load balancing across pods
- Health-check based routing
- Session affinity support

**Ingress Controller**:
- Layer 7 load balancing
- SSL termination
- Path-based routing
- Rate limiting support

## Security Implementation

### Application Security

#### Session Management
```go
// HTTP-only cookies
http.SetCookie(w, &http.Cookie{
    Name:     "session_token",
    Value:    sessionToken,
    Expires:  time.Now().Add(sessionTTL),
    Path:     "/",
    HttpOnly: true,  // Not accessible via JavaScript
    SameSite: http.SameSiteLaxMode,
})
```

**Features**:
- Secure session tokens (32-byte random)
- HTTP-only cookies (XSS protection)
- Session TTL with automatic cleanup
- MongoDB TTL indexes for expired sessions

#### Password Security
```go
// BCrypt hashing
hashedPassword, _ := bcrypt.GenerateFromPassword(
    []byte(password), 
    bcrypt.DefaultCost,
)
```

**Features**:
- BCrypt password hashing
- Salt per password
- Configurable work factor
- Resistant to rainbow table attacks

### Container Security

#### Multi-stage Builds
```dockerfile
# Builder stage
FROM golang:1.25-alpine AS builder
# Build application

# Final stage
FROM alpine:latest
# Copy only binary
```

**Benefits**:
- Minimal attack surface
- No build tools in production
- Smaller image size
- Faster deployments

#### Non-root User
```dockerfile
# Run as non-root user
USER nobody
```

#### Read-only Root Filesystem
```yaml
securityContext:
  readOnlyRootFilesystem: true
  runAsNonRoot: true
  runAsUser: 65534
```

### Network Security

#### Network Policies
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: todo-app-policy
spec:
  podSelector:
    matchLabels:
      app: todo-app
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: ingress-nginx
    ports:
    - protocol: TCP
      port: 8080
```

#### TLS/SSL Configuration
```yaml
annotations:
  cert-manager.io/cluster-issuer: "letsencrypt-prod"
  nginx.ingress.kubernetes.io/ssl-redirect: "true"
```

### Secrets Management

**Kubernetes Secrets**:
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: mongodb-secret
type: Opaque
data:
  username: <base64-encoded>
  password: <base64-encoded>
```

**External Secrets** (Recommended for production):
- AWS Secrets Manager
- Azure Key Vault
- HashiCorp Vault
- Google Secret Manager

## Resource Management

### Resource Limits

**Application Resources**:
```yaml
resources:
  requests:
    cpu: "500m"      # Guaranteed CPU
    memory: "512Mi"  # Guaranteed memory
  limits:
    cpu: "2000m"     # Maximum CPU
    memory: "2048Mi" # Maximum memory
```

**Why This Matters**:
- Prevents resource exhaustion
- Ensures fair scheduling
- Enables horizontal scaling
- Protects cluster stability

### Horizontal Pod Autoscaling

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: todo-app-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: todo-app
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

**Benefits**:
- Automatic scaling based on load
- Cost optimization
- Better resource utilization
- Handles traffic spikes

### Quality of Service

**Guaranteed QoS**:
```yaml
resources:
  requests:
    cpu: "1000m"
    memory: "1024Mi"
  limits:
    cpu: "1000m"
    memory: "1024Mi"
```
- Highest priority
- Never evicted for resources
- Suitable for critical workloads

## Persistent Storage

### Volume Management

**PersistentVolumeClaim**:
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: todo-app-uploads-pvc
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: fast-ssd
  resources:
    requests:
      storage: 10Gi
```

**Features**:
- Data persistence across pod restarts
- Automatic provisioning
- Storage class selection
- Backup capabilities

### Backup Strategy

**Volume Snapshots**:
```yaml
apiVersion: snapshot.storage.k8s.io/v1
kind: VolumeSnapshot
metadata:
  name: todo-app-snapshot
spec:
  volumeSnapshotClassName: csi-snapshot-class
  source:
    persistentVolumeClaimName: todo-app-uploads-pvc
```

**MongoDB Backups**:
```bash
# MongoDB dump
kubectl exec deployment/mongodb -n todo-app -- \
  mongodump --out=/backup --gzip

# Copy from pod
kubectl cp todo-app/mongodb-pod:/backup ./backup
```

## Monitoring and Observability

### Metrics Collection

**Prometheus Integration**:
```yaml
apiVersion: v1
kind: Service
metadata:
  name: todo-app
  annotations:
    prometheus.io/scrape: "true"
    prometheus.io/port: "8080"
    prometheus.io/path: "/metrics"
```

**Key Metrics**:
- Request rate
- Error rate
- Response time
- Resource usage
- Database connections

### Logging

**Structured Logging**:
```go
log.Printf("[INFO] User %s logged in at %s", 
    username, time.Now())
```

**Log Aggregation**:
- ELK Stack (Elasticsearch, Logstash, Kibana)
- Loki + Grafana
- AWS CloudWatch
- Azure Monitor

### Distributed Tracing

**OpenTelemetry**:
```go
import "go.opentelemetry.io/otel"

func handler(w http.ResponseWriter, r *http.Request) {
    ctx, span := tracer.Start(r.Context(), "handler")
    defer span.End()
    // ... handler code
}
```

## Deployment Strategies

### Rolling Updates (Default)

```yaml
strategy:
  type: RollingUpdate
  rollingUpdate:
    maxSurge: 1        # Extra pod during update
    maxUnavailable: 0  # Zero downtime
```

**Benefits**:
- Zero downtime
- Gradual rollout
- Easy rollback
- Production-safe

### Blue-Green Deployment

```bash
# Deploy new version (green)
kubectl apply -f deployment-green.yaml

# Test green deployment
kubectl port-forward svc/todo-app-green 8080:8080

# Switch traffic to green
kubectl patch service todo-app -p \
  '{"spec":{"selector":{"version":"green"}}}'

# Remove blue deployment
kubectl delete deployment todo-app-blue
```

### Canary Deployment

```yaml
# 10% of traffic to canary
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: todo-app
spec:
  http:
  - match:
    - headers:
        canary:
          exact: "true"
    route:
    - destination:
        host: todo-app-canary
  - route:
    - destination:
        host: todo-app-stable
      weight: 90
    - destination:
        host: todo-app-canary
      weight: 10
```

## Health Checks

### Liveness Probe

```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 30
  periodSeconds: 10
  timeoutSeconds: 5
  failureThreshold: 3
```

**Purpose**: Restart unhealthy pods

### Readiness Probe

```yaml
readinessProbe:
  httpGet:
    path: /ready
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 5
  timeoutSeconds: 3
  successThreshold: 1
```

**Purpose**: Remove unready pods from load balancing

### Startup Probe

```yaml
startupProbe:
  httpGet:
    path: /startup
    port: 8080
  initialDelaySeconds: 0
  periodSeconds: 10
  failureThreshold: 30
```

**Purpose**: Allow slow-starting containers

## GitOps Benefits

### Declarative Configuration

**Everything as Code**:
- Infrastructure configuration
- Application configuration
- Deployment manifests
- Network policies

**Benefits**:
- Version controlled
- Auditable changes
- Easy rollback
- Reproducible deployments

### Automated Sync

**ArgoCD Features**:
- Automatic deployment on Git changes
- Self-healing (corrects manual changes)
- Pruning (removes deleted resources)
- Diff visualization

### Pull-based Deployment

**Security Advantages**:
- No cluster credentials in CI
- Cluster pulls changes
- Reduced attack surface
- Better secret management

## Disaster Recovery

### Backup Strategy

**Application Data**:
- PVC snapshots
- MongoDB dumps
- Configuration backups
- Secret backups

**Recovery Procedures**:
```bash
# Restore from snapshot
kubectl apply -f volume-snapshot-restore.yaml

# Restore MongoDB
kubectl exec -it deployment/mongodb -n todo-app -- \
  mongorestore --gzip /backup
```

### High Availability Architecture

```
┌─────────────────────────────────────────┐
│  Multi-Zone Kubernetes Cluster          │
├─────────────────────────────────────────┤
│  Zone A    │   Zone B   │   Zone C     │
│  App Pod 1 │  App Pod 2 │  App Pod 3   │
│  MongoDB 1 │  MongoDB 2 │  MongoDB 3   │
│  (Primary) │  (Secondary)│ (Secondary)  │
└─────────────────────────────────────────┘
```

## Performance Optimization

### Caching Strategy

**Application-level**:
- Session caching
- Query result caching
- Static asset caching

**Infrastructure-level**:
- CDN for static assets
- Database connection pooling
- Image caching in registry

### Resource Optimization

**Efficient Images**:
- Multi-stage builds
- Minimal base images (Alpine)
- Layer optimization
- Build caching

**Efficient Queries**:
- Database indexes
- Query optimization
- Connection pooling
- Prepared statements

## Compliance and Auditing

### Audit Logging

**Kubernetes Audit**:
```yaml
apiVersion: audit.k8s.io/v1
kind: Policy
rules:
  - level: Metadata
    resources:
    - group: ""
      resources: ["pods"]
```

**Application Audit**:
- User actions
- Authentication attempts
- Data modifications
- System events

### Compliance Features

- **RBAC**: Role-based access control
- **Network Policies**: Traffic restrictions
- **Pod Security Standards**: Security constraints
- **Secret Encryption**: At-rest encryption
- **Audit Logs**: Complete audit trail

## Next Steps

- [Review configuration options](configuration.md)
- [Learn monitoring techniques](monitoring.md)
- [Understand learning outcomes](learning-outcomes.md)

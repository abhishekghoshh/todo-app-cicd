# Monitoring and Troubleshooting

Comprehensive guide for monitoring, debugging, and troubleshooting the Todo App CI/CD pipeline.

## Monitoring Tools

### ArgoCD Dashboard

Access the ArgoCD UI for visual monitoring:

```bash
# Port forward to ArgoCD
kubectl port-forward svc/argocd-server -n argocd 8080:443

# Get admin password
kubectl -n argocd get secret argocd-initial-admin-secret \
  -o jsonpath="{.data.password}" | base64 -d
```

**Access**: `http://localhost:8080`

**Features**:
- Application sync status
- Resource health
- Deployment history
- Rollback capabilities
- Live resource view

### Kubernetes Dashboard

```bash
# Install Kubernetes Dashboard
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.7.0/aio/deploy/recommended.yaml

# Create admin user
kubectl create serviceaccount dashboard-admin -n kubernetes-dashboard
kubectl create clusterrolebinding dashboard-admin \
  --clusterrole=cluster-admin \
  --serviceaccount=kubernetes-dashboard:dashboard-admin

# Get token
kubectl -n kubernetes-dashboard create token dashboard-admin

# Port forward
kubectl port-forward -n kubernetes-dashboard service/kubernetes-dashboard 8443:443
```

**Access**: `https://localhost:8443`

## Monitoring Commands

### Application Status

```bash
# Check ArgoCD applications
kubectl get applications -n argocd

# Detailed application status
argocd app get todo-app

# Application sync status
argocd app list
```

### Pod Monitoring

```bash
# List all pods
kubectl get pods -n todo-app

# Watch pods in real-time
kubectl get pods -n todo-app -w

# Check pod details
kubectl describe pod <pod-name> -n todo-app

# Pod resource usage
kubectl top pods -n todo-app
```

### Service Status

```bash
# List services
kubectl get svc -n todo-app

# Service endpoints
kubectl get endpoints -n todo-app

# Describe service
kubectl describe svc todo-app-svc -n todo-app
```

### Deployment Status

```bash
# Check deployments
kubectl get deployments -n todo-app

# Deployment details
kubectl describe deployment todo-app -n todo-app

# Rollout status
kubectl rollout status deployment/todo-app -n todo-app

# Rollout history
kubectl rollout history deployment/todo-app -n todo-app
```

### Ingress Monitoring

```bash
# Check ingress
kubectl get ingress -n todo-app

# Ingress details
kubectl describe ingress todo-app-ingress -n todo-app

# Check ingress controller
kubectl get pods -n ingress-nginx
```

### Storage Monitoring

```bash
# Check persistent volumes
kubectl get pv

# Check persistent volume claims
kubectl get pvc -n todo-app

# PVC details
kubectl describe pvc todo-app-uploads-pvc -n todo-app
```

## Log Management

### Application Logs

```bash
# Stream application logs
kubectl logs -f deployment/todo-app -n todo-app

# Last 100 lines
kubectl logs --tail=100 deployment/todo-app -n todo-app

# Logs from specific container
kubectl logs -f deployment/todo-app -c app -n todo-app

# Previous pod logs (after crash)
kubectl logs deployment/todo-app --previous -n todo-app
```

### MongoDB Logs

```bash
# Stream MongoDB logs
kubectl logs -f deployment/mongodb -n todo-app

# Check MongoDB connection issues
kubectl logs deployment/mongodb -n todo-app | grep -i "connection"
```

### ArgoCD Logs

```bash
# ArgoCD application controller
kubectl logs -f deployment/argocd-application-controller -n argocd

# ArgoCD server logs
kubectl logs -f deployment/argocd-server -n argocd

# ArgoCD repo server (for sync issues)
kubectl logs -f deployment/argocd-repo-server -n argocd
```

### GitHub Actions Logs

Navigate to: `https://github.com/<owner>/<repo>/actions`

**Check for**:
- Build failures
- Push failures
- Workflow execution times
- Error messages

## Common Issues and Solutions

### 1. Application Not Syncing

**Symptoms**:
- ArgoCD shows "OutOfSync"
- Changes not reflected in cluster

**Diagnosis**:
```bash
# Check sync status
argocd app get todo-app

# Check for sync errors
kubectl describe application todo-app -n argocd
```

**Solutions**:
```bash
# Force sync
argocd app sync todo-app --force

# Refresh application
argocd app refresh todo-app

# Check ArgoCD repo server logs
kubectl logs -f deployment/argocd-repo-server -n argocd
```

### 2. Pods Crashing (CrashLoopBackOff)

**Symptoms**:
- Pods repeatedly restarting
- `CrashLoopBackOff` status

**Diagnosis**:
```bash
# Check pod status
kubectl get pods -n todo-app

# View pod events
kubectl describe pod <pod-name> -n todo-app

# Check logs
kubectl logs <pod-name> -n todo-app
kubectl logs <pod-name> --previous -n todo-app
```

**Common Causes**:
1. **MongoDB connection failure**
   ```bash
   # Check MongoDB status
   kubectl get pods -l app=mongodb -n todo-app
   
   # Test connection
   kubectl exec -it <app-pod> -n todo-app -- \
     nc -zv mongodb-svc 27017
   ```

2. **Missing environment variables**
   ```bash
   # Check environment
   kubectl exec -it <app-pod> -n todo-app -- env | grep MONGO
   ```

3. **Resource limits**
   ```bash
   # Check resource usage
   kubectl top pod <pod-name> -n todo-app
   ```

### 3. ImagePullBackOff

**Symptoms**:
- Pods stuck in `ImagePullBackOff`
- Cannot pull container image

**Diagnosis**:
```bash
# Check pod events
kubectl describe pod <pod-name> -n todo-app | grep -A 10 Events
```

**Solutions**:
1. **Check image exists**
   ```bash
   docker pull abhishek1009/todo-app:latest
   ```

2. **Verify image tag**
   ```bash
   # Check deployment image
   kubectl get deployment todo-app -n todo-app -o yaml | grep image:
   ```

3. **Check pull secrets** (if private registry)
   ```bash
   kubectl get secrets -n todo-app
   ```

### 4. Ingress Not Working

**Symptoms**:
- Cannot access application via URL
- 404 or 502 errors

**Diagnosis**:
```bash
# Check ingress
kubectl get ingress -n todo-app
kubectl describe ingress todo-app-ingress -n todo-app

# Check ingress controller
kubectl get pods -n ingress-nginx
kubectl logs -f deployment/ingress-nginx-controller -n ingress-nginx
```

**Solutions**:
1. **Verify ingress class**
   ```bash
   kubectl get ingressclass
   ```

2. **Check service endpoints**
   ```bash
   kubectl get endpoints todo-app-svc -n todo-app
   ```

3. **Test service directly**
   ```bash
   kubectl port-forward svc/todo-app-svc -n todo-app 3000:8080
   curl http://localhost:3000
   ```

### 5. Persistent Volume Issues

**Symptoms**:
- PVC stuck in `Pending`
- Data not persisting

**Diagnosis**:
```bash
# Check PVC status
kubectl get pvc -n todo-app

# Check PV
kubectl get pv

# Describe PVC
kubectl describe pvc <pvc-name> -n todo-app
```

**Solutions**:
1. **Check storage class**
   ```bash
   kubectl get storageclass
   ```

2. **Verify volume provisioner**
   ```bash
   kubectl get pods -n kube-system | grep provisioner
   ```

3. **Check node capacity**
   ```bash
   kubectl describe nodes | grep -A 5 "Allocated resources"
   ```

### 6. GitHub Actions Failing

**Symptoms**:
- Workflow fails
- Image/chart not published

**Diagnosis**:
- Check workflow logs in GitHub Actions tab
- Review error messages

**Common Issues**:

1. **Docker Hub authentication**
   ```yaml
   # Verify secrets are set
   Settings → Secrets → Actions
   - DOCKERHUB_TOKEN
   ```

2. **Build failures**
   ```bash
   # Test build locally
   docker build -t test .
   ```

3. **Permission issues**
   ```yaml
   # Check workflow permissions
   permissions:
     contents: write
     pull-requests: write
   ```

### 7. Helm Chart Issues

**Symptoms**:
- Chart installation fails
- Values not applied

**Diagnosis**:
```bash
# Test Helm template
helm template todo-app ./helm \
  -f argocd/apps/todo-app/overlays/dev/values.yaml

# Check Helm releases
helm list -n todo-app

# Helm release status
helm status todo-app -n todo-app
```

**Solutions**:
```bash
# Validate chart
helm lint ./helm

# Dry-run installation
helm install todo-app ./helm --dry-run --debug

# Check chart values
helm get values todo-app -n todo-app
```

## Performance Monitoring

### Resource Usage

```bash
# Node resource usage
kubectl top nodes

# Pod resource usage
kubectl top pods -n todo-app

# Detailed metrics
kubectl describe node <node-name> | grep -A 5 "Allocated resources"
```

### Application Metrics

```bash
# Port forward to application
kubectl port-forward svc/todo-app-svc -n todo-app 8080:8080

# Check application health (if endpoint exists)
curl http://localhost:8080/health

# Load testing
ab -n 1000 -c 10 http://localhost:8080/
```

### Database Performance

```bash
# Connect to MongoDB
kubectl exec -it deployment/mongodb -n todo-app -- mongosh

# MongoDB stats
db.serverStatus()
db.stats()

# Check connections
db.serverStatus().connections
```

## Health Checks

### Application Health

```bash
# Check if pods are ready
kubectl get pods -n todo-app

# Readiness probe status
kubectl describe pod <pod-name> -n todo-app | grep -A 5 Readiness

# Liveness probe status
kubectl describe pod <pod-name> -n todo-app | grep -A 5 Liveness
```

### Cluster Health

```bash
# Check node status
kubectl get nodes

# Check system pods
kubectl get pods -n kube-system

# Check component status
kubectl get componentstatuses
```

## Debugging Tools

### Interactive Shell

```bash
# Access application container
kubectl exec -it deployment/todo-app -n todo-app -- /bin/sh

# Access MongoDB container
kubectl exec -it deployment/mongodb -n todo-app -- mongosh
```

### Network Debugging

```bash
# Test DNS resolution
kubectl run -it --rm debug --image=busybox --restart=Never -- nslookup mongodb-svc.todo-app.svc.cluster.local

# Test connectivity
kubectl run -it --rm debug --image=busybox --restart=Never -- wget -O- http://todo-app-svc:8080

# Network policies
kubectl get networkpolicies -n todo-app
```

### Resource Debugging

```bash
# Check resource quotas
kubectl get resourcequota -n todo-app

# Check limit ranges
kubectl get limitrange -n todo-app

# Node resources
kubectl describe node <node-name>
```

## Rollback Procedures

### ArgoCD Rollback

```bash
# View deployment history
argocd app history todo-app

# Rollback to previous version
argocd app rollback todo-app <revision-number>

# Sync to specific revision
argocd app sync todo-app --revision <revision>
```

### Kubernetes Rollback

```bash
# Rollout history
kubectl rollout history deployment/todo-app -n todo-app

# Rollback to previous version
kubectl rollout undo deployment/todo-app -n todo-app

# Rollback to specific revision
kubectl rollout undo deployment/todo-app --to-revision=2 -n todo-app
```

### Helm Rollback

```bash
# Helm release history
helm history todo-app -n todo-app

# Rollback to previous release
helm rollback todo-app -n todo-app

# Rollback to specific revision
helm rollback todo-app 2 -n todo-app
```

## Cleanup and Maintenance

### Clean Up Failed Pods

```bash
# Delete failed pods
kubectl delete pods --field-selector status.phase=Failed -n todo-app

# Delete evicted pods
kubectl delete pods --field-selector status.phase=Evicted -n todo-app
```

### Clean Up Old Resources

```bash
# Delete completed jobs
kubectl delete jobs --field-selector status.successful=1 -n todo-app

# Clean up old replicasets
kubectl delete replicaset --all -n todo-app
```

### Restart Deployments

```bash
# Restart application
kubectl rollout restart deployment/todo-app -n todo-app

# Restart MongoDB
kubectl rollout restart deployment/mongodb -n todo-app

# Restart all deployments
kubectl rollout restart deployment --all -n todo-app
```

## Monitoring Best Practices

1. **Regular Health Checks**: Monitor application and database health daily
2. **Log Aggregation**: Consider using ELK stack or Loki for centralized logging
3. **Metrics Collection**: Set up Prometheus and Grafana for detailed metrics
4. **Alerting**: Configure alerts for critical issues
5. **Backup Strategy**: Regular backups of persistent data
6. **Documentation**: Keep runbooks updated with common issues

## Next Steps

- [Learn about production features](production-features.md)
- [Review configuration options](configuration.md)
- [Understand CI/CD workflows](cicd-workflows.md)

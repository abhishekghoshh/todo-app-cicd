# todo-app-cicd
A full CICD for a TODO app


## Helm Package

> I found this article [How to Build & Push Helm Chart to Docker Hub](https://forums.docker.com/t/how-to-build-push-helm-chart-to-docker-hub/133548)


## ArgoCD installation

**Configure ArgoCD in local**
```zsh
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```

#!/usr/bin/env bash
set -e

echo "➡️ Enabling Helm support inside Kustomize for ArgoCD..."

# Ensure namespace exists
kubectl create namespace argocd 2>/dev/null || true

# Patch argocd-cm to add build options
echo "📌 Patching argocd-cm with --enable-helm"

kubectl patch configmap argocd-cm -n argocd \
  --type merge \
  -p '{"data":{"kustomize.buildOptions":"--enable-helm --load-restrictor LoadRestrictionsNone"}}' || {

  echo "❗ argocd-cm not found. Creating a new one..."

  kubectl apply -n argocd -f - <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-cm
  namespace: argocd
data:
  kustomize.buildOptions: "--enable-helm --load-restrictor LoadRestrictionsNone"
EOF
}

echo "🔄 Restarting ArgoCD components so changes take effect..."

kubectl rollout restart deployment argocd-repo-server -n argocd
kubectl rollout restart deployment argocd-server -n argocd
kubectl rollout restart deployment argocd-applicationset-controller -n argocd
kubectl rollout restart deployment argocd-dex-server -n argocd 2>/dev/null || true

echo "⏳ Waiting for components to become ready..."

kubectl wait deployment/argocd-repo-server \
  -n argocd --for=condition=Available=True --timeout=120s

kubectl wait deployment/argocd-server \
  -n argocd --for=condition=Available=True --timeout=120s

echo "🎉 Helm support in Kustomize is now ENABLED in ArgoCD!"

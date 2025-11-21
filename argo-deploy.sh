#!/usr/bin/env bash
set -e

if [ -z "$1" ]; then
  echo "❌ Please provide an environment: dev | stage | prod"
  echo "Or provide destroy to uninstall ArgoCD"
  exit 1
fi

if [ "$1" == "destroy" ]; then
    echo "🗑️ Uninstalling ArgoCD..."
    kubectl delete -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml || echo "⚠️ ArgoCD resources not found, skipping deletion"
    kubectl delete namespace argocd || echo "⚠️ ArgoCD namespace not found, skipping deletion"
    echo "✅ ArgoCD uninstalled."
    exit 0
fi

TARGET_ENV="$1"
ROOT_FILE="application.yaml"
TEMP_FILE="application-rendered.yaml"

echo "➡️ Using environment: $TARGET_ENV"

# ---------------------------------------
# Step 1: Install / Start ArgoCD
# ---------------------------------------
echo "📌 Checking if ArgoCD is installed..."

if kubectl get namespace argocd >/dev/null 2>&1 && kubectl get deployment argocd-server -n argocd >/dev/null 2>&1; then
    echo "✅ ArgoCD is already installed, skipping installation"
else
    echo "📌 Installing ArgoCD..."
    
    kubectl create namespace argocd 2>/dev/null || true
    kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
    
    echo "⏳ Waiting for ArgoCD to be ready..."
    
    kubectl wait deployment/argocd-server \
        -n argocd \
        --for=condition=Available=True \
        --timeout=180s
    
    echo "✅ ArgoCD is ready!"
fi

# ---------------------------------------
# Step 2: Substitute the environment
# ---------------------------------------
echo "🔧 Rendering root application..."

# Substitute ${targetEnvironment} → actual env
TARGET_ENV="$TARGET_ENV" envsubst < "$ROOT_FILE" > "$TEMP_FILE"

echo "📄 Rendered file:"
cat "$TEMP_FILE"

# ---------------------------------------
# Step 3: Apply the file
# ---------------------------------------
echo "🚀 Applying root application..."
kubectl apply -f "$TEMP_FILE"

# delete the temp file
rm "$TEMP_FILE"

echo "✅ Applied root application."

echo "🎉 Successfully deployed for environment: $TARGET_ENV"

# Argocd initial admin password
echo "ℹ️ ArgoCD initial admin password:"
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d; echo
echo "🔗 Access ArgoCD UI at: http://localhost:8080 (port-forwarding)"
kubectl port-forward svc/argocd-server -n argocd 8080:443 >/dev/null 2>&1 &
echo "🚀 Port-forwarding started..."
echo "Use username: admin"
echo "Use the above password to login."
echo "Press Ctrl+C to stop port-forwarding."
# Keep the script running to maintain port-forwarding
wait
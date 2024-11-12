#!/usr/bin/env bash

set -o errexit -o nounset -o pipefail

echo "Deploying a sample Pod..."
echo
kubectl apply -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: outyet
  labels:
    app: outyet
spec:
  securityContext:
    runAsNonRoot: true
  containers:
    - name: outyet
      image: ghcr.io/l0rd/outyet:latest
      ports:
        - containerPort: 8080
          protocol: TCP
      resources:
        limits:
          memory: "128Mi"
          cpu: "500m"
      securityContext:
        allowPrivilegeEscalation: false
        capabilities:
          drop:
            - "ALL"
        seccompProfile:
          type: RuntimeDefault
---
kind: Service
apiVersion: v1
metadata:
  labels:
    app: outyet
  name: outyet-service
spec:
  ports:
  - port: 8080
    protocol: TCP
    targetPort: 8080
  selector:
    app: outyet
  type: ClusterIP
---
apiVersion: route.openshift.io/v1
kind: Route
metadata:
  labels:
    app: outyet
  name: outyet-route
spec:
  port:
    targetPort: 8080
  to:
    kind: Service
    name: outyet-service
---
EOF

kubectl wait --for=condition=Ready pod/outyet

echo
echo "Installing..."
echo
go install cmd/kubectl-debug_ide.go

echo
echo "Running..."
echo
kubectl debug-ide outyet \
  --image ghcr.io/l0rd/outyet-dev:latest \
  --copy-to outyet-debug \
  --share-processes \
  --git-repository https://github.com/l0rd/outyet.git

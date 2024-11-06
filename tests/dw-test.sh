#!/usr/bin/env bash

set -o errexit -o nounset -o pipefail

kubectl apply -f - <<EOF
apiVersion: workspace.devfile.io/v1alpha2
kind: DevWorkspace
metadata:
  name: outyet-dw
  namespace: mloriedo-dev
  annotations:
    controller.devfile.io/debug-start: "true"
spec:
  contributions:
    - components:
        - container:
            env:
              - name: CODE_HOST
                value: 0.0.0.0
          name: che-code-runtime-description
      name: che-code
      uri: https://eclipse-che.github.io/che-plugin-registry/main/v3/plugins/che-incubator/che-code/latest/devfile.yaml
  started: true
  template:
    metadata:
      name: outyet-dw
    attributes:
      controller.devfile.io/storage-type: ephemeral
      pod-overrides:
        spec:
          shareProcessNamespace: true
    components:
      - container:
          image: ghcr.io/l0rd/outyet-dev:latest
          memoryRequest: 2G
          memoryLimit: 10G
          cpuRequest: '1'
          cpuLimit: '4'
          mountSources: true
        name: dev
      - container:
          cpuLimit: 500m
          endpoints:
            - exposure: public
              name: port8080
              path: /
              protocol: http
              secure: false
              targetPort: 8080
          image: ghcr.io/l0rd/outyet:latest
          memoryLimit: 128Mi
          mountSources: false
        name: outyet
    projects:
      - git:
          remotes:
            origin: https://github.com/l0rd/outyet.git
        name: outyet
        sourceType: Git
EOF

#POD=$(kubectl get pod -l "controller.devfile.io/devworkspace_name"="outyet-dw" -o jsonpath="{ ..metadata.name}")
#export POD
##kubectl wait --for=condition=ready --timeout=20s pod/"${POD}"
#while true
#do
#  sleep 1
#  kubectl exec -it "${POD}" -c cde -- /bin/sh -c "cat /checode/entrypoint-logs.txt" || true
#done

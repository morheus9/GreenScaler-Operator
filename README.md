![Logo](assets/logo.jpg)

## For developer

CRD Example:

```Yaml
apiVersion: app.example.com/v1alpha1
kind: GreenScalerService
metadata:
  labels:
    app.kubernetes.io/name: greenscaler-operator
    app.kubernetes.io/managed-by: kustomize
  name: greenscalerservice-sample
spec:
  timeZone: "UTC"
  targets:
    - kind: Deployment
      name: my-app
      namespace: default
  schedule:
    - from: "09:00"
      to: "18:00"
      replicas: 3
    - from: "18:00"
      to: "09:00"
      replicas: 0
```

### 1. Install utilities

- kubebuilder

```
go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest

kubebuilder init --domain morheus9.dev --repo github.com/morheus9/GreenScaler-Operator
kubebuilder create api --group greenscaler --version v1alpha1 --kind GreenScalerService --resource --controller
```

- controller-gen

```
make controller-gen
GOBIN=/home/pi/go/bin go install sigs.k8s.io/controller-tools/cmd/controller-gen@latest
```

- kustomize

```
make kustomize
sudo apt install kustomize
```

### 2. Code generation

##### Generates code based on markers in types.go

```
make generate
```

##### Generates CRD, RBAC, Webhook manifests

```
make manifests
```

### 3. Build

```
make build
```

### 4. Installing in cluster

```
make install
```

### 5. Start in development mode

```
make run
```

##### Checking

```
kubectl get crd | grep
```

PS

```
make generate
make manifests
make build

make install
make deploy IMG=morheus/GreenScaler-operator:0.0.1
kubectl apply -k config/samples/

make uninstall
make undeploy IMG=morheus/GreenScaler-operator:0.0.1
kubectl apply -k config/samples/
```

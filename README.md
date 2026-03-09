## Auto observability operator

---

## For developer

### 1. Install utilities

- kubebuilder

```
go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest

kubebuilder init --domain example.com --repo github.com/morheus9/GreenScaler-Operator
kubebuilder create api --group app --version v1alpha1 --kind ServiceMonitoring --resource --controller
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

##### Генерирует код на основе маркеров в types.go

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
make deploy IMG=morheus/observability-operator:0.0.1
kubectl apply -k config/samples/

make uninstall
make undeploy IMG=morheus/observability-operator:0.0.1
kubectl apply -k config/samples/
```

```
┌──────────────────────────────────────────────────────────────────────────┐
│                        Kubernetes Cluster                                │
│                                                                          │
│  ┌────────────────────────────────────────────────────────────────────┐  │
│  │                     ArgoCD Applications                            │  │
│  │  ┌─────────┐ ┌──────────┐ ┌─────────┐ ┌─────────────┐ ┌─────────┐  │  │
│  │  │ PromOps │ │JaegerOps │ │FluentOps│ │Observability│ │ Grafana │  │  │
│  │  │ App     │ │ App      │ │ App     │ │ Operator    │ │ App     │  │  │
│  │  └────┬────┘ └────┬─────┘ └────┬────┘ └────┬────────┘ └────┬────┘  │  │
│  └───────┼───────────┼────────────┼───────────┼───────────────┼───────┘  │
│          │           │            │           │               │          │
│  ┌───────▼───────────▼────────────▼───────────▼───────────────▼───────┐  │
│  │                      Операторы (Operators)                         │  │
│  │  ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐       │  │
│  │  │ Prometheus      │ │Jaeger           │ │ Fluent          │       │  │
│  │  │ Operator        │ │Operator         │ │ Operator        │       │  │
│  │  │ • ServiceMonitor│ │• Jaeger CRD     │ │ • LogPipeline   │       │  │
│  │  │ • PrometheusRule│ │• Instrumentation│ │ • LogParser     │       │  │
│  │  └──────┬──────────┘ └────────┬────────┘ └───────┬─────────┘       │  │
│  │         │                     │                  │                 │  │
│  │  ┌──────▼─────────────────────▼──────────────────▼──────────────┐  │  │
│  │  │                    Инструменты (Tools)                       │  │  │
│  │  │  ┌──────────┐    ┌────────────┐  ┌────────────────────────┐  │  │  │
│  │  │  │Prometheus│    │  Jaeger    │  │      Fluent Bit        │  │  │  │
│  │  │  │ • Сбор   │    │ • Collector│  │ • DaemonSet на нодах   │  │  │  │
│  │  │  │   метрик │    │ • Query    │  │ • Сбор логов           │  │  │  │
│  │  │  │ • Storage│    │ • UI       │  │ • Отправка в Loki      │  │  │  │
│  │  │  └─────┬────┘    └─────┬──────┘  └──────────┬─────────────┘  │  │  │
│  │  │        │               │                    │                │  │  │
│  │  │  ┌─────▼───────────────▼────────────────────▼─────────────┐  │  │  │
│  │  │  │                    Хранилища (Storage)                 │  │  │  │
│  │  │  │  ┌────────────┐   ┌──────────┐   ┌──────────────────┐  │  │  │  │
│  │  │  │  │ TSDB       │   │Cassandra │   │      Loki        │  │  │  │  │
│  │  │  │  │(Prometheus)│   │/Elastic  │   │ • Хранение логов │  │  │  │  │
│  │  │  │  │            │   │(Jaeger)  │   │ • LogQL          │  │  │  │  │
│  │  │  │  └────────────┘   └──────────┘   └──────────────────┘  │  │  │  │
│  │  │  └────────────────────────────────────────────────────────┘  │  │  │
│  │  └──────────────────────────────────────────────────────────────┘  │  │
│  │                                                                    │  │
│  │  ┌──────────────────────────────────────────────────────────────┐  │  │
│  │  │             Пользовательские приложения                      │  │  │
│  │  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐           │  │  │
│  │  │  │   App 1     │  │   App 2     │  │   App 3     │           │  │  │
│  │  │  │ • Метрики   │  │ • Метрики   │  │ • Метрики   │           │  │  │
│  │  │  │ • Логи      │  │ • Логи      │  │ • Логи      │           │  │  │
│  │  │  │ • Трейсы    │  │ • Трейсы    │  │ • Трейсы    │           │  │  │
│  │  │  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘           │  │  │
│  │  │         │                │                │                  │  │  │
│  │  │  ┌──────▼────────────────▼────────────────▼───────────────┐  │  │  │
│  │  │  │             Observability-Operator                     │  │  │  │
│  │  │  │                                                        │  │  │  │
│  │  │  │ • Отслеживает Deployments/Services                     │  │  │  │
│  │  │  │ • Создает CRD для других операторов                    │  │  │  │
│  │  │  │ • Настраивает Observability автоматически              │  │  │  │
│  │  │  └────────────────────────────────────────────────────────┘  │  │  │
│  │  └──────────────────────────────────────────────────────────────┘  │  │
│  └────────────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────────────┘
```

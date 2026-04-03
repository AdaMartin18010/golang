# Helm Charts Design

> **分类**: 工程与云原生
> **标签**: #helm #kubernetes #charts #packaging #templating
> **参考**: Helm Best Practices, Chart Museum, Helm Hub

---

## 1. Formal Definition

### 1.1 What is Helm?

Helm is a package manager for Kubernetes that simplifies deployment, versioning, and management of applications. Helm Charts are packages of pre-configured Kubernetes resources that define, install, and upgrade complex Kubernetes applications.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Helm Architecture                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                        HELM CLIENT                                 │   │
│   │                                                                      │   │
│   │  helm install     helm upgrade    helm rollback    helm uninstall   │   │
│   │       │                │               │               │            │   │
│   │       └────────────────┴───────────────┴───────────────┘            │   │
│   │                          │                                          │   │
│   │                          ▼                                          │   │
│   │              ┌─────────────────────┐                                │   │
│   │              │   TILLER/HELM 3     │  (In-cluster or client-only)  │   │
│   │              │   (Release Manager) │                                │   │
│   │              └──────────┬──────────┘                                │   │
│   │                         │                                           │   │
│   └─────────────────────────┼───────────────────────────────────────────┘   │
│                             │                                               │
│                             ▼                                               │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                     HELM CHART STRUCTURE                             │   │
│   │                                                                      │   │
│   │  mychart/                                                            │   │
│   │  ├── Chart.yaml           # Chart metadata                           │   │
│   │  ├── values.yaml          # Default configuration values             │   │
│   │  ├── values.schema.json   # JSON schema for validation               │   │
│   │  ├── charts/              # Sub-charts dependencies                  │   │
│   │  ├── templates/           # Kubernetes manifest templates            │   │
│   │  │   ├── _helpers.tpl     # Named templates/definitions              │   │
│   │  │   ├── deployment.yaml  # Deployment template                      │   │
│   │  │   ├── service.yaml     # Service template                         │   │
│   │  │   ├── ingress.yaml     # Ingress template                         │   │
│   │  │   ├── configmap.yaml   # ConfigMap template                       │   │
│   │  │   ├── secret.yaml      # Secret template                          │   │
│   │  │   ├── hpa.yaml         # Horizontal Pod Autoscaler                │   │
│   │  │   ├── pdb.yaml         # Pod Disruption Budget                    │   │
│   │  │   ├── sa.yaml          # Service Account                          │   │
│   │  │   ├── NOTES.txt        # Post-installation notes                  │   │
│   │  │   └── tests/           # Test hooks                               │   │
│   │  └── README.md            # Documentation                            │   │
│   │                                                                      │   │
│   │  Templating Flow:                                                   │   │
│   │  values.yaml + templates/*.yaml ──► Rendering ──► K8s manifests      │   │
│   │                                                                      │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│   RELEASE MANAGEMENT:                                                       │
│   • Revision history stored in ConfigMaps/Secrets                           │
│   • Atomic releases (all-or-nothing)                                        │
│   • Automatic rollback capability                                           │
│   • Upgrade hooks for migrations                                            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Helm Chart Lifecycle

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Helm Chart Lifecycle                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   DEVELOPMENT                  TESTING                     DISTRIBUTION     │
│      │                            │                            │            │
│      ▼                            ▼                            ▼            │
│   ┌───────┐                   ┌───────┐                   ┌──────────┐      │
│   │Create │                   │Template│                  │  Package │      │
│   │Chart  │──────► lint ─────►│Render │──────► test ────►│   Push   │      │
│   │       │                   │       │                   │   to     │      │
│   │• Chart│                   │• Values│                  │ Registry │      │
│   │  yaml │                   │  validation                │          │      │
│   │• Templates│               │• Template                  │• Version │      │
│   │• Values│                  │  syntax                    │  tagging │      │
│   └───────┘                   └───────┘                   └──────────┘      │
│                                                                             │
│   DEPLOYMENT                   OPERATIONS                    CLEANUP        │
│      │                            │                            │            │
│      ▼                            ▼                            ▼            │
│   ┌─────────┐                 ┌─────────┐                 ┌──────────┐      │
│   │ Install │                 │ Upgrade │                 │Uninstall │      │
│   │  /      │◄─── rollback ──│  /      │                 │  /       │      │
│   │ Upgrade │                 │ Rollback│                 │ Delete   │      │
│   └────┬────┘                 └────┬────┘                 └──────────┘      │
│        │                           │                                        │
│        └──► Release created ──────►└──► Release updated                      │
│                   │                         │                                │
│                   ▼                         ▼                                │
│            ┌─────────────────────────────────────┐                           │
│            │         RELEASE SECRETS/            │                           │
│            │         CONFIGMAPS                  │                           │
│            │  (Stores release history, values,   │                           │
│            │   and rendered templates)           │                           │
│            └─────────────────────────────────────┘                           │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Implementation Patterns

### 2.1 Production-Ready Chart Structure

```yaml
# Chart.yaml
apiVersion: v2
name: my-application
description: A production-ready Helm chart for My Application
type: application
version: 1.2.3
appVersion: "2.0.0"
kubeVersion: ">=1.21.0-0"

keywords:
  - myapp
  - backend
  - api

home: https://github.com/example/my-application
sources:
  - https://github.com/example/my-application

maintainers:
  - name: John Doe
    email: john.doe@example.com
    url: https://github.com/johndoe

dependencies:
  - name: postgresql
    version: 12.x.x
    repository: https://charts.bitnami.com/bitnami
    condition: postgresql.enabled
    tags:
      - database

  - name: redis
    version: 17.x.x
    repository: https://charts.bitnami.com/bitnami
    condition: redis.enabled

  - name: common
    version: 2.x.x
    repository: https://charts.bitnami.com/bitnami
    tags:
      - bitnami-common

annotations:
  category: ApplicationServer
  licenses: Apache-2.0
  images: |
    - name: my-application
      image: docker.io/example/my-application:2.0.0
```

```yaml
# values.yaml - Production defaults
# Global settings
global:
  imageRegistry: ""
  imagePullSecrets: []
  storageClass: ""

# Image configuration
image:
  registry: docker.io
  repository: example/my-application
  tag: ""
  pullPolicy: IfNotPresent
  pullSecrets: []
  debug: false

# Common labels
commonLabels: {}
commonAnnotations: {}

# Replica configuration
replicaCount: 3
minReplicas: 2
maxReplicas: 10

# Deployment strategy
updateStrategy:
  type: RollingUpdate
  rollingUpdate:
    maxSurge: 25%
    maxUnavailable: 0

# Pod Disruption Budget
pdb:
  enabled: true
  minAvailable: 2
  maxUnavailable: ""

# Horizontal Pod Autoscaler
autoscaling:
  enabled: true
  minReplicas: 3
  maxReplicas: 20
  targetCPUUtilizationPercentage: 70
  targetMemoryUtilizationPercentage: 80
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
        - type: Percent
          value: 10
          periodSeconds: 60
    scaleUp:
      stabilizationWindowSeconds: 0
      policies:
        - type: Percent
          value: 100
          periodSeconds: 15
        - type: Pods
          value: 4
          periodSeconds: 15
      selectPolicy: Max

# Service configuration
service:
  type: ClusterIP
  port: 8080
  targetPort: http
  nodePort: ""
  annotations: {}

# Ingress configuration
ingress:
  enabled: true
  className: nginx
  annotations:
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/proxy-body-size: "10m"
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
  hosts:
    - host: api.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: api-example-com-tls
      hosts:
        - api.example.com

# Service Account
serviceAccount:
  create: true
  annotations: {}
  name: ""
  automountServiceAccountToken: false

# Security Context
podSecurityContext:
  runAsNonRoot: true
  runAsUser: 1000
  runAsGroup: 1000
  fsGroup: 1000

containerSecurityContext:
  allowPrivilegeEscalation: false
  readOnlyRootFilesystem: true
  capabilities:
    drop:
      - ALL
  seccompProfile:
    type: RuntimeDefault

# Resource limits
resources:
  limits:
    cpu: 1000m
    memory: 1Gi
  requests:
    cpu: 200m
    memory: 256Mi

# Health checks
livenessProbe:
  enabled: true
  httpGet:
    path: /health/live
    port: http
  initialDelaySeconds: 30
  periodSeconds: 10
  timeoutSeconds: 5
  failureThreshold: 3
  successThreshold: 1

readinessProbe:
  enabled: true
  httpGet:
    path: /health/ready
    port: http
  initialDelaySeconds: 10
  periodSeconds: 5
  timeoutSeconds: 3
  failureThreshold: 3
  successThreshold: 1

startupProbe:
  enabled: true
  httpGet:
    path: /health/ready
    port: http
  initialDelaySeconds: 10
  periodSeconds: 5
  timeoutSeconds: 3
  failureThreshold: 30

# Persistence
persistence:
  enabled: false
  storageClass: ""
  accessMode: ReadWriteOnce
  size: 10Gi
  existingClaim: ""

# Configuration
config:
  LOG_LEVEL: info
  PORT: "8080"
  METRICS_PORT: "9090"

# Secrets (references)
existingSecret: ""

# Environment variables
extraEnvVars: []
#  - name: CUSTOM_VAR
#    value: custom_value

# ConfigMap files
extraConfigMaps: {}
#  config.json: |
#    {"key": "value"}

# Extra volumes
extraVolumes: []
#  - name: extra-volume
#    emptyDir: {}

extraVolumeMounts: []
#  - name: extra-volume
#    mountPath: /extra

# Init containers
initContainers: []
#  - name: init-db
#    image: busybox:latest
#    command: ['sh', '-c', 'echo init']

# Sidecar containers
sidecars: []

# Pod topology spread constraints
topologySpreadConstraints: []
#  - maxSkew: 1
#    topologyKey: topology.kubernetes.io/zone
#    whenUnsatisfiable: ScheduleAnyway
#    labelSelector:
#      matchLabels:
#        app.kubernetes.io/name: my-application

# Affinity
affinity:
  podAntiAffinity:
    preferredDuringSchedulingIgnoredDuringExecution:
      - weight: 100
        podAffinityTerm:
          labelSelector:
            matchExpressions:
              - key: app.kubernetes.io/name
                operator: In
                values:
                  - my-application
          topologyKey: kubernetes.io/hostname

# Tolerations
tolerations: []

# Node selector
nodeSelector: {}

# Pod annotations
podAnnotations: {}

# Network Policy
networkPolicy:
  enabled: true
  ingress:
    - from:
        - namespaceSelector:
            matchLabels:
              name: ingress-nginx
  egress:
    - to:
        - namespaceSelector: {}

# Metrics
metrics:
  enabled: true
  service:
    port: 9090
  serviceMonitor:
    enabled: true
    namespace: ""
    interval: 30s
    scrapeTimeout: 10s
    labels: {}
    relabelings: []
    metricRelabelings: []
  prometheusRule:
    enabled: false
    rules: []

# Tracing
tracing:
  enabled: false
  exporter: otlp
  endpoint: ""
  samplingRate: 0.1

# Feature flags
features:
  caching: true
  rateLimiting: true
  requestLogging: true
```

### 2.2 Template Best Practices

```yaml
# templates/_helpers.tpl
{{/* vim: set filetype=mustache: */}}

{{/* Expand the name of the chart */}}
{{- define "myapp.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/* Create a default fully qualified app name */}}
{{- define "myapp.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/* Create chart name and version */}}
{{- define "myapp.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/* Common labels */}}
{{- define "myapp.labels" -}}
helm.sh/chart: {{ include "myapp.chart" . }}
{{ include "myapp.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- if .Values.commonLabels }}
{{ toYaml .Values.commonLabels }}
{{- end }}
{{- end }}

{{/* Selector labels */}}
{{- define "myapp.selectorLabels" -}}
app.kubernetes.io/name: {{ include "myapp.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/* Create the name of the service account */}}
{{- define "myapp.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "myapp.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/* Validate required values */}}
{{- define "myapp.validateValues" -}}
{{- if not .Values.image.repository -}}
myapp: image.repository
    The image repository is required.
{{- end -}}
{{- end -}}

{{/* Image definition */}}
{{- define "myapp.image" -}}
{{- $registryName := .Values.image.registry -}}
{{- $repositoryName := .Values.image.repository -}}
{{- $tag := .Values.image.tag | default .Chart.AppVersion -}}
{{- if $registryName }}
{{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
{{- else -}}
{{- printf "%s:%s" $repositoryName $tag -}}
{{- end -}}
{{- end -}}
```

```yaml
# templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "myapp.fullname" . }}
  labels:
    {{- include "myapp.labels" . | nindent 4 }}
  {{- with .Values.commonAnnotations }}
  annotations:
    {{- toYaml . | nindent 4 }}
  {{- end }}
spec:
  {{- if not .Values.autoscaling.enabled }}
  replicas: {{ .Values.replicaCount }}
  {{- end }}
  {{- if .Values.updateStrategy }}
  strategy:
    {{- toYaml .Values.updateStrategy | nindent 4 }}
  {{- end }}
  selector:
    matchLabels:
      {{- include "myapp.selectorLabels" . | nindent 6 }}
  template:
    metadata:
      annotations:
        checksum/config: {{ include (print $.Template.BasePath "/configmap.yaml") . | sha256sum }}
        checksum/secrets: {{ include (print $.Template.BasePath "/secret.yaml") . | sha256sum }}
        {{- with .Values.podAnnotations }}
        {{- toYaml . | nindent 8 }}
        {{- end }}
      labels:
        {{- include "myapp.labels" . | nindent 8 }}
    spec:
      {{- include "myapp.imagePullSecrets" . | nindent 6 }}
      serviceAccountName: {{ include "myapp.serviceAccountName" . }}
      automountServiceAccountToken: {{ .Values.serviceAccount.automountServiceAccountToken }}
      securityContext:
        {{- toYaml .Values.podSecurityContext | nindent 8 }}

      {{- if .Values.initContainers }}
      initContainers:
        {{- toYaml .Values.initContainers | nindent 8 }}
      {{- end }}

      containers:
        - name: {{ .Chart.Name }}
          securityContext:
            {{- toYaml .Values.containerSecurityContext | nindent 12 }}
          image: {{ include "myapp.image" . }}
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          ports:
            - name: http
              containerPort: {{ .Values.service.targetPort }}
              protocol: TCP
            {{- if .Values.metrics.enabled }}
            - name: metrics
              containerPort: {{ .Values.metrics.service.port }}
              protocol: TCP
            {{- end }}
          env:
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: POD_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
            {{- range $key, $value := .Values.config }}
            - name: {{ $key }}
              value: {{ $value | quote }}
            {{- end }}
            {{- if .Values.existingSecret }}
            - name: DATABASE_URL
              valueFrom:
                secretKeyRef:
                  name: {{ .Values.existingSecret }}
                  key: database-url
            {{- end }}
            {{- with .Values.extraEnvVars }}
            {{- toYaml . | nindent 12 }}
            {{- end }}
          envFrom:
            - configMapRef:
                name: {{ include "myapp.fullname" . }}
          {{- if .Values.livenessProbe.enabled }}
          livenessProbe:
            {{- toYaml (omit .Values.livenessProbe "enabled") | nindent 12 }}
          {{- end }}
          {{- if .Values.readinessProbe.enabled }}
          readinessProbe:
            {{- toYaml (omit .Values.readinessProbe "enabled") | nindent 12 }}
          {{- end }}
          {{- if .Values.startupProbe.enabled }}
          startupProbe:
            {{- toYaml (omit .Values.startupProbe "enabled") | nindent 12 }}
          {{- end }}
          resources:
            {{- toYaml .Values.resources | nindent 12 }}
          volumeMounts:
            - name: tmp
              mountPath: /tmp
            {{- if .Values.persistence.enabled }}
            - name: data
              mountPath: /data
            {{- end }}
            {{- with .Values.extraVolumeMounts }}
            {{- toYaml . | nindent 12 }}
            {{- end }}
        {{- with .Values.sidecars }}
        {{- toYaml . | nindent 8 }}
        {{- end }}

      volumes:
        - name: tmp
          emptyDir: {}
        {{- if .Values.persistence.enabled }}
        - name: data
          {{- if .Values.persistence.existingClaim }}
          persistentVolumeClaim:
            claimName: {{ .Values.persistence.existingClaim }}
          {{- else }}
          persistentVolumeClaim:
            claimName: {{ include "myapp.fullname" . }}
          {{- end }}
        {{- end }}
        {{- with .Values.extraVolumes }}
        {{- toYaml . | nindent 8 }}
        {{- end }}

      {{- with .Values.nodeSelector }}
      nodeSelector:
        {{- toYaml . | nindent 8 }}
      {{- end }}

      {{- with .Values.affinity }}
      affinity:
        {{- toYaml . | nindent 8 }}
      {{- end }}

      {{- with .Values.topologySpreadConstraints }}
      topologySpreadConstraints:
        {{- toYaml . | nindent 8 }}
      {{- end }}

      {{- with .Values.tolerations }}
      tolerations:
        {{- toYaml . | nindent 8 }}
      {{- end }}
```

---

## 3. Production-Ready Configurations

### 3.1 Multi-Environment Values

```yaml
# values-production.yaml
# Production-specific overrides

replicaCount: 5
minReplicas: 3
maxReplicas: 50

resources:
  limits:
    cpu: 2000m
    memory: 2Gi
  requests:
    cpu: 500m
    memory: 512Mi

autoscaling:
  enabled: true
  minReplicas: 5
  maxReplicas: 50
  targetCPUUtilizationPercentage: 60
  targetMemoryUtilizationPercentage: 70

pdb:
  enabled: true
  minAvailable: 3

networkPolicy:
  enabled: true
  ingress:
    - from:
        - namespaceSelector:
            matchLabels:
              name: ingress-nginx
        - podSelector:
            matchLabels:
              app.kubernetes.io/name: monitoring
  egress:
    - to:
        - namespaceSelector:
            matchLabels:
              name: database
    - to:
        - namespaceSelector:
            matchLabels:
              name: cache

config:
  LOG_LEVEL: warn
  CACHE_ENABLED: "true"
  RATE_LIMIT_ENABLED: "true"
```

---

## 4. Security Considerations

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       Helm Security Best Practices                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  CHART DEVELOPMENT                                                          │
│  ✓ Sign charts with GPG or cosign                                           │
│  ✓ Scan images for vulnerabilities                                          │
│  ✓ No secrets in values.yaml (use sealed-secrets or external)               │
│  ✓ Use read-only root filesystem                                            │
│  ✓ Drop all capabilities                                                    │
│  ✓ Run as non-root user                                                     │
│                                                                             │
│  DEPLOYMENT                                                                 │
│  ✓ Verify chart signatures before install                                   │
│  ✓ Use specific chart versions (no latest)                                  │
│  ✓ Review rendered templates before apply                                   │
│  ✓ Use values.schema.json for validation                                    │
│  ✓ Enable audit logging                                                     │
│                                                                             │
│  SECRETS MANAGEMENT                                                         │
│  ✓ External secret operators (Vault, AWS Secrets Manager)                   │
│  ✓ Sealed Secrets for GitOps                                                │
│  ✓ SOPS for encrypted values files                                          │
│  ✓ Never commit plain secrets to Git                                        │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 5. Decision Matrices

### 5.1 Helm vs Alternatives

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Deployment Tool Comparison                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Tool           │  Use Case                    │  When to Choose            │
├─────────────────┼──────────────────────────────┼────────────────────────────│
│  Helm           │  Package management          │  Multiple apps, versioning │
│                 │  Templating                  │  Community charts          │
│  ───────────────┼──────────────────────────────┼────────────────────────────│
│  Kustomize      │  Native K8s patching         │  GitOps, overlays          │
│                 │  No templating               │  Simple configurations     │
│  ───────────────┼──────────────────────────────┼────────────────────────────│
│  Operator       │  Complex stateful apps       │  DBs, MQs, caches          │
│                 │  Lifecycle management        │  Day-2 operations          │
│  ───────────────┼──────────────────────────────┼────────────────────────────│
│  Terraform      │  Infrastructure + K8s        │  Multi-cloud, IaC          │
│                 │  External resources          │  State management          │
│  ───────────────┼──────────────────────────────┼────────────────────────────│
│  CD (Argo/Flux) │  GitOps delivery             │  Continuous deployment     │
│                 │  Progressive delivery        │  Multi-cluster             │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 6. Best Practices Summary

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Helm Best Practices Summary                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  CHART DESIGN                                                               │
│  ✓ Follow semantic versioning                                               │
│  ✓ Use helper templates for consistency                                     │
│  ✓ Provide sensible defaults in values.yaml                                 │
│  ✓ Document all values in README                                            │
│  ✓ Include values.schema.json for validation                                │
│                                                                             │
│  TEMPLATING                                                                 │
│  ✓ Use nindent for proper indentation                                       │
│  ✓ Quote all string values                                                  │
│  ✓ Check for empty values before rendering                                  │
│  ✓ Use include instead of template                                          │
│  ✓ Keep templates DRY                                                       │
│                                                                             │
│  TESTING                                                                    │
│  ✓ Use helm lint for syntax checking                                        │
│  ✓ Template rendering tests                                                 │
│  ✓ Chart testing (ct) for integration tests                                 │
│  ✓ Test upgrades from previous versions                                     │
│                                                                             │
│  DISTRIBUTION                                                               │
│  ✓ Sign charts before publishing                                            │
│  ✓ Use chart repositories (Harbor, ChartMuseum)                             │
│  ✓ Version charts independently of app                                      │
│  ✓ Maintain changelog                                                       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## References

1. Helm Best Practices
2. Chart Testing Tools
3. Helm Security Guide
4. Artifact Hub
5. Helm Charts Repository

/*
 * Copyright (C) 2019 Rohith Jayawardene <gambol99@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package bootstrap

const (
	// NamespaceAdminClusterRole is used by the namespace admin
	NamespaceAdminClusterRole = `
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: hub:system:ns-admin
rules:
- nonResourceURLs:
  - /swagger*
  - /swaggerapi
  - /swaggerapi/*
  - /version
  verbs:
  - get
- apiGroups:
  - apps
  - batch
  - extensions
  - networking.k8s.io
  resources:
  - cronjobs
  - deployments
  - deployments/rollback
  - deployments/scale
  - ingresses
  - jobs
  - networkpolicies
  - replicasets
  - replicasets/scale
  - replicationcontrollers/scale
  - statefulsets
  - statefulsets/scale
  verbs:
  - '*'
- apiGroups:
  - policy
  resources:
  - poddisruptionbudgets
  verbs:
  - '*'
- apiGroups:
  - ""
  resources:
  - configmaps
  - endpoints
  - persistentvolumeclaims
  - persistentvolumes
  - pods
  - pods/attach
  - pods/exec
  - pods/log
  - pods/portforward
  - secrets
  - serviceaccounts
  - services
  verbs:
  - '*'
- apiGroups:
  - autoscaling
  resources:
  - "*"
  verbs:
  - "*"
- apiGroups:
  - '*'
  resources:
  - '*'
  verbs:
  - get
  - watch
  - list
- apiGroups:
  - certmanager.k8s.io
  resources:
  - certificates
  - challenges
  - orders
  verbs:
  - "*"
`

	// BootstrapJobTemplate is the template for the job
	BootstrapJobTemplate = `
---
apiVersion: batch/v1
kind: Job
metadata:
  name: bootstrap
  namespace: kube-system
spec:
  backoffLimit: 20
  template:
    spec:
      serviceAccountName: "hub-admin"
      restartPolicy: OnFailure
      containers:
        - name: bootstrap
          image: {{ .BootImage }}
          imagePullPolicy: Always
          env:
            - name: CONFIG_DIR
              value: "/config"
            - name: PROVIDER
              value: "{{ .Provider }}"
            - name: OLM_VERSION
              value: "{{ .OLMVersion }}"
          volumeMounts:
            - name: bundle
              mountPath: /config/bundles
            - name: olm
              mountPath: /config/olm
      volumes:
        - name: bundle
          configMap:
            name: bootstrap
        - name: olm
          configMap:
            name: bootstrap-olm
`

	// BootstrapJobConfigmap is the configmap
	BootstrapJobConfigmap = `
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: bootstrap
  namespace: kube-system
data:
  repositories: |
    # @TODO need to figure out a way to install this via the OLM - it's
    # due the use of a genCert which is causing the OLM to cycle the
    # deployments
    svc-cat,https://svc-catalog-charts.storage.googleapis.com
  charts: |
    # helm source for the service catalog
    svc-cat/catalog,catalog,--values /config/bundles/catalog.yaml
    # The values supplied to the service-catalog
  catalog.yaml: |
    imagePullPolicy: IfNotPresent
    apiserver:
      storage:
        etcd:
          image: quay.io/coreos/etcd:v3.4.1@sha256:49d3d4a81e0d030d3f689e7167f23e120abf955f7d08dbedf3ea246485acee9f
          imagePullPolicy: IfNotPresent
          persistence:
            enabled: true
            size: 4Gi
    controllerManager:
      annotations:
        prometheus.io/scheme: https
      brokerRelistInterval: 20m
      enablePrometheusScrape: true
      resyncInterval: 5m
`

	// BootstrapJobOLMConfig is the configuration for the olm
	BootstrapJobOLMConfig = `
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: bootstrap-olm
  namespace: kube-system
data:
  ## Metrics namespace
  namespaces.yaml: |
  {{- range .Namespaces }}
    ---
    apiVersion: v1
    kind: Namespace
    metadata:
      name: {{ .Name }}
      {{- if .EnableIstio }}
      labels:
        'istio-injection': 'enabled'
      {{- end }}
    ---
    apiVersion: operators.coreos.com/v1
    kind: OperatorGroup
    metadata:
      name: operator-group
      namespace: {{ .Name }}
    spec:
      targetNamespaces:
        - {{ .Name }}
  {{- end }}
  ## Catalog
  catalog.yaml: |
    ---
    apiVersion: operators.coreos.com/v1alpha1
    kind: CatalogSource
    metadata:
      name: appvia-catalog
      namespace: olm
    spec:
      {{- if .Catalog.Image }}
      image: quay.io/appvia/operator-catalog:{{ .Catalog.Image }}
      {{- else }}
      addr:  {{ .Catalog.GRPC }}
      {{- end }}
      displayName: Appvia Operators
      publisher: Appvia.io
      sourceType: grpc
  ## Operator groups
  {{- if .OperatorGroups }}
  operatorgroups.yaml: |
    {{- range .OperatorGroups }}
    ---
    apiVersion: operators.coreos.com/v1
    kind: OperatorGroup
    metadata:
      name: operator-group
      namespace: {{ . }}
    spec:
      namespaces:
        - {{ . }}
    {{- end }}
  {{- end }}
  ## Subscriptions
  {{- range $i, $x := .Operators }}
  subscription-{{ $i }}.yaml: |
    # operator_selector: {{ $x.Label }}
    apiVersion: operators.coreos.com/v1alpha1
    kind: Subscription
    metadata:
      name: {{ $x.Package }}
      namespace: {{ or $x.Namespace "prometheus" }}
    spec:
      name: {{ $x.Package }}
      channel: {{ $x.Channel }}
      installPlanApproval: {{ or $x.InstallPlan "Automatic" }}
      source: {{ or $x.Catalog "appvia-catalog" }}
      sourceNamespace: olm
  {{- end }}
  ## External DNS
  crd-external-dns.yaml: |
    ---
    apiVersion: helm.appvia.io/v1alpha1
    kind: ExternalDns
    metadata:
      name: external-dns
      namespace: kube-dns
    spec:
      {{- if eq .Provider "eks" }}
      provider: aws
      aws:
        credentials:
          accessKey: {{ .Credentials.AWS.AccessKey }}
          mountPath: "/.aws"
          secretKey: {{ .Credentials.AWS.SecretKey }}
          region: {{ .Credentials.AWS.Region }}
      {{- end }}
      {{- if eq .Provider "gke" }}
      provider: google
      google:
        # @TODO need to change this to one with reduced perms
        serviceAccountKey: '{{ .Credentials.GKE.Account | toJson }}'
      {{- end }}
      domainFilters:
        - {{ .Domain }}
      policy: sync
      metrics:
        enabled: true
      rbac:
        create: true
        serviceAccountName: external-dns
      service:
        annotations:
          prometheus.io/scrape: 'true'
          prometheus.io/port: '7979'
      sources:
        - ingress
        - service
      fullnameOverride: external-dns
  crd-monitoring.yaml: |
    ---
    apiVersion: helm.appvia.io/v1alpha1
    kind: Metrics
    metadata:
      name: metrics
      namespace: prometheus
    spec:
      alertmanager:
        enabled: true
        alertmanagerSpec:
          image:
            repository: quay.io/prometheus/alertmanager
        tag: v0.19.0
        replicas: 1
        retention: 120h
      coreDns:
        enabled: true
        service:
          port: 10054
          targetPort: 10054
          selector:
            k8s-app: kube-dns
      kubeProxy:
        enabled: true
      kubeStateMetrics:
        enabled: true
      kubelet:
        enabled: true
      nodeExporter:
        enabled: true
      prometheus:
        prometheusSpec:
          image:
            repository: quay.io/prometheus/prometheus
            tag: v2.12.0
          replicas: 2
          retention: 10d
          ruleSelector:
            app: prometheus
          serviceMonitorSelector:
            metrics: prometheus
      kube-state-metrics:
        fullnameOverride: kube-state-metrics
      prometheus-node-exporter:
        fullnameOverride: node-exporter
  ## Cloud Service Brokers
  {{ if and (eq .Provider "eks") ( .EnableServiceBroker ) }}
  crd-aws-service-broker.yaml: |
    apiVersion: helm.appvia.io/v1alpha1
    kind: AwsServicebroker
    metadata:
      name: aws-broker
      namespace: brokers
    spec:
      image: awsservicebroker/aws-servicebroker:beta
      aws:
        accesskeyid: {{ .Credentials.AWS.AccessKey }}
        bucket: awsservicebroker
        key: templates/latest
        region: {{ .Credentials.AWS.Region }}
        s3region: us-east-1
        secretkey: {{ .Credentials.AWS.SecretKey }}
        tablename: awssb
        targetaccountid: {{ .Credentials.AWS.AccountID }}
  {{- end }}
  {{ if and (eq .Provider "gke") ( .EnableServiceBroker ) }}
  crd-gcp-service-broker.yaml: |
    {{- if .EnableIstio }}
    # Due to the MySQL protocol we need to mitigate this when in
    # permissive mode
    # https://istio.io/faq/security/#mysql-with-mtls
    ---
    apiVersion: "authentication.istio.io/v1alpha1"
    kind: Policy
    metadata:
      name: gcp-broker-db-mtls
      namespace: brokers
    spec:
      targets:
        - name: gcp-broker-db
    {{- end }}
    ---
    apiVersion: helm.appvia.io/v1alpha1
    kind: Mariadb
    metadata:
      name: gcp-broker-db
      namespace: brokers
    spec:
      db:
        forcePassword: false
        name: {{ .Broker.Name }}
      master:
        persistence:
          enabled: true
          size: 4Gi
      rootUser:
        forcePassword: true
        password: {{ .Broker.Database.Password }}
      fullnameOverride: gcp-broker-db
      metrics:
        enabled: true
        serviceMonitor:
          enabled: false
      replication:
        enabled: false
      serviceAccount:
        create: true
      slave:
        replicas: 0
    ---
    apiVersion: helm.appvia.io/v1alpha1
    kind: GcpServiceBroker
    metadata:
      name: gcp-broker
      namespace: brokers
    spec:
      broker:
        password: {{ .Broker.Password }}
        service_account_json: '{{ .Credentials.GKE.Account | toJson }}'
        username: {{ .Broker.Username }}
      image:
        repository: gcr.io/gcp-service-broker/gcp-service-broker
        tag: v4.3.0
      mysql:
        embedded: false
        host: gcp-broker-db
        mysqlDatabase: {{ .Broker.Database.Name }}
        mysqlPassword: {{ .Broker.Database.Password }}
        mysqlUser: root
      replicaCount: 1
  {{- end }}
  {{ if .EnableKiali }}
  crd-kiali.yaml: |
    ---
    apiVersion: v1
    kind: Secret
    metadata:
      name: kiali
      namespace: istio-system
    type: Opaque
    data:
      passphrase: {{ .Kiali.Password | b64enc }}
      username: {{ "admin" | b64enc }}
    ---
    apiVersion: kiali.io/v1alpha1
    kind: Kiali
    metadata:
      name: kiali
      namespace: istio-system
    spec:
      installation_tag: Appvia
      istio_namespace: istio-system
      deployment:
        namespace: istio-system
        verbose_mode: '4'
        view_only_mode: false
      external_services:
        grafana:
          url: 'http://grafana-service.grafana.svc.cluster.local:3000'
        prometheus:
          url: 'http://prometheus.prometheus.svc.cluster.local:9090'
        #tracing:
        #  url: ''
      server:
        web_root: "/kiali"
  {{- end }}
`
)

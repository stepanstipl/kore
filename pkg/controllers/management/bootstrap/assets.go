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
    # svc-cat/catalog,catalog,--values /config/bundles/catalog.yaml
    # The values supplied to the service-catalog
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
        serviceAccountSecret: google
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

  ## Grafana
  crd-grafana.yaml: |
    ---
    apiVersion: helm.appvia.io/v1alpha1
    kind: Mariadb
    metadata:
      name: grafana-db
      namespace: grafana
    spec:
      db:
        forcePassword: false
        name: grafana
      master:
        persistence:
          enabled: true
          size: 10Gi
      rootUser:
        forcePassword: true
        password: {{ .Grafana.Database.Password }}
      fullnameOverride: grafana-db
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
    apiVersion: integreatly.org/v1alpha1
    kind: Grafana
    metadata:
      name: grafana
      namespace: grafana
    spec:
      initialReplicas: 3
      {{- if eq .Provider "eks" }} 
      service:
        type: LoadBalancer
        annotations:
          'service.beta.kubernetes.io/aws-load-balancer-backend-protocol': 'http'
          'external-dns.alpha.kubernetes.io/hostname': '{{ .Grafana.Hostname }}'
      {{- else }}
      ingress:
        enabled: true
        hostname: {{ .Grafana.Hostname }}
      service:
        type: NodePort
      {{- end }}
      config:
        analytics:
          check_for_updates: true
        auth:
          disable_signout_menu: false
        auth.basic:
          enabled: false
        auth.anonymous:
          enabled: true
        auth.generic_oauth:
          allow_sign_up: true
          enabled: true
          client_id : {{ .Grafana.ClientID }}
          client_secret: {{ .Grafana.ClientSecret }}
          scopes: email,profile
          api_url: {{ .Grafana.UserInfoURL }}
          auth_url: {{ .Grafana.AuthURL }}
          token_url: {{ .Grafana.TokenURL }}
          #allowed_domains: {{ .Domain }}
        database:
          host: grafana-db
          name: grafana
          password: {{ .Grafana.Database.Password }}
          type: mysql
          user: root
        log:
          level: info
          mode: console
        paths:
          data: /var/lib/grafana/data
          logs: /var/log/grafana
          plugins: /var/lib/grafana/plugins
          provisioning: /etc/grafana/provisioning
        security:
          admin_password: {{ .Grafana.Password }}
          admin_user: admin
        server:
          domain: {{ .Grafana.Hostname }}
          enable_gzip: true
          root_url: http://{{ .Grafana.Hostname }}
        users:
          auto_assign_org_role: Editor
      dashboardLabelSelector:
        - matchExpressions:
          - key: app
            operator: In
            values:
              - grafana

  ## Logging
  crd-logging.yaml: |
    ---
    apiVersion: helm.appvia.io/v1alpha1
    kind: Loki
    metadata:
      name: loki
      namespace: logging
    spec:
      loki:
        enabled: true
        image:
          repository: grafana/loki
          tag: v0.3.0
        persistence:
          accessModes:
            - ReadWriteOnce
          enabled: true
          size: 10Gi
          storageClassName: {{ .StorageClass }}
        replicas: 1
        serviceMonitor:
          enabled: true
          additionalLabels:
            metrics: prometheus

      promtail:
        enabled: true
        image:
          pullPolicy: IfNotPresent
          repository: grafana/promtail
          tag: v0.3.0
        serviceMonitor:
          additionalLabels:
            metrics: prometheus
    ---
    apiVersion: integreatly.org/v1alpha1
    kind: GrafanaDataSource
    metadata:
      name: loki
      namespace: grafana
    spec:
      name: logging.yaml
      datasources:
        - access: proxy
          editable: false
          isDefault: false
          name: loki
          type: loki
          url: http://loki.logging.svc.cluster.local:3100
          version: 1
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
`
)

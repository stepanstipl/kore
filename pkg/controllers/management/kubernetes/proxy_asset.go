/**
 * Copyright 2020 Appvia Ltd <info@appvia.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package kubernetes

// AuthProxyDeployment is the deployment template for the authentication proxy
const AuthProxyDeployment = `
apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    kore.appvia.io/owner: "true"
  name: proxy
  namespace: {{ .Namespace }}
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: kore:oidc:proxy
rules:
  - apiGroups:
      - '*'
    resources:
      - users
      - groups
      - serviceaccount
    verbs:
      - impersonate
  - apiGroups:
      - authentication.k8s.io
    resources:
      - userextras/scopes
      - tokenreviews
    verbs:
      - create
      - impersonate
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: kore:oidc:proxy
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: kore:oidc:proxy
subjects:
  - kind: ServiceAccount
    name: proxy
    namespace: kore
---
apiVersion: v1
kind: Service
metadata:
  annotations:
    external-dns.alpha.kubernetes.io/hostname: api.{{ .Domain }}
    {{- if eq .Provider "EKS" }}
    service.beta.kubernetes.io/aws-load-balancer-proxy-protocol: '*'
    {{- end }}
  name: proxy
  namespace: {{ .Namespace }}
spec:
  externalTrafficPolicy: Local
  ports:
  - name: https
    port: 443
    protocol: TCP
    targetPort: 10443
  selector:
    name: {{ .Deployment }}
  sessionAffinity: None
  type: LoadBalancer
{{- if and .TLSKey .TLSCert }}
---
apiVersion: v1
kind: Secret
metadata:
  name: tls
  namespace: {{ .Namespace }}
  annotations:
    kore.appvia.io/owned: "true"
data:
  tls.crt: {{ toString .TLSCert | b64enc }}
  tls.key: {{ toString .TLSKey | b64enc }}
{{- end }}
---
apiVersion: v1
kind: Secret
metadata:
  name: ca
  namespace: {{ .Namespace }}
  annotations:
    kore.appvia.io/owned: "true"
data:
  ca.crt: {{ toString .CACert | b64enc }}
---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    kore.appvia.io/owned: "true"
  labels:
    name: {{ .Deployment }}
  name: {{ .Deployment }}
  namespace: {{ .Namespace }}
spec:
  replicas: {{ .Replicas }}
  selector:
    matchLabels:
      name: {{ .Deployment }}
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      labels:
        name: {{ .Deployment }}
      annotations:
        prometheus.io/port: "8080"
        prometheus.io/scrape: "true"
    spec:
      serviceAccountName: proxy
      containers:
        - name: oidc-proxy
          image: {{ .Image }}
          env:
            {{- if .OpenID }}
            - name: IDP_CLIENT_ID
              value: {{ .ClientID }}
            - name: IDP_SERVER_URL
              value: {{ .ServerURL }}
            {{- end }}
            - name: JWT_SIGNER_CERT
              value: /ca/ca.crt
            - name: TLS_CERT
              value: /tls/tls.crt
            - name: TLS_KEY
              value: /tls/tls.key
          args:
            {{- range .AllowedIPs }}
            - --allowed-ips={{ . }}
            {{- end }}
            {{- if .OpenID }}
            {{- range .UserClaims }}
            - --idp-user-claims={{ . }}
            {{- end }}
            {{- end }}
            - --verifiers=localjwt
            {{- if .OpenID }}
            - --verifiers=openid
            {{- end }}
            - --verifiers=tokenreview
          ports:
            - containerPort: 10443
              protocol: TCP
            - containerPort: 8080
              protocol: TCP
          readinessProbe:
            failureThreshold: 3
            periodSeconds: 10
            successThreshold: 1
            timeoutSeconds: 1
            {{- if eq .Provider "EKS" }}
            tcpSocket:
              port: 10443
            {{- else }}
            httpGet:
              path: /ready
              port: 10443
              scheme: HTTPS
            {{- end }}
          resources:
            resources:
              limits:
                cpu: 50m
                memory: "64Mi"
              requests:
                cpu: 10m
                memory: "12Mi"
          volumeMounts:
            - mountPath: /tls
              name: tls
              readOnly: true
            - mountPath: /ca
              name: ca
              readOnly: true
      volumes:
        - name: ca
          secret:
            defaultMode: 420
            secretName: ca
        - name: tls
          secret:
            defaultMode: 420
            secretName: tls
---
apiVersion: autoscaling/v2beta2
kind: HorizontalPodAutoscaler
metadata:
  annotations:
    kore.appvia.io/owned: "true"
  name: {{ .Deployment }}
  namespace: {{ .Namespace }}
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: {{ .Deployment }}
  minReplicas: {{ .Replicas }}
  maxReplicas: {{ .MaxReplicas }}
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 50
`

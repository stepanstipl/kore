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

package assets

// LocalCompose is used by the local start/stop commands to run kore locally in Docker Compose.
const LocalCompose = `
---
version: '3'
services:
  etcd:
    image: bitnami/etcd:3.4.4
    environment:
      ALLOW_NONE_AUTHENTICATION: "yes"
    ports:
      - 2379:2379

  kube-controller-manager:
    image: gcr.io/google-containers/kube-controller-manager-amd64:v1.15.11
    command:
      - /usr/local/bin/kube-controller-manager
      - --master=http://kube-apiserver:8080

  kube-apiserver:
    image: gcr.io/google-containers/kube-apiserver-amd64:v1.15.11
    command:
      - /usr/local/bin/kube-apiserver
      - --address=0.0.0.0
      - --alsologtostderr
      - --authorization-mode=RBAC
      - --bind-address=0.0.0.0
      - --default-watch-cache-size=200
      - --delete-collection-workers=10
      - --etcd-servers=http://etcd:2379
      - --log-flush-frequency=10s
      - --runtime-config=autoscaling/v1=false
      - --runtime-config=autoscaling/v2beta1=false
      - --runtime-config=autoscaling/v2beta2=false
      - --runtime-config=batch/v1=false
      - --runtime-config=batch/v1beta1=false
      - --runtime-config=networking.k8s.io/v1=false
      - --runtime-config=networking.k8s.io/v1beta1=false
      - --runtime-config=node.k8s.io/v1beta1=false
    ports:
      - 8080:8080
      - 6443:6443

  database:
    image: mariadb:10.5.1
    environment:
      MYSQL_ROOT_PASSWORD: pass
    entrypoint:
      sh -c "
        echo 'CREATE DATABASE IF NOT EXISTS kore;' > /docker-entrypoint-initdb.d/init.sql;
        /usr/local/bin/docker-entrypoint.sh --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci"
    ports:
      - 3306:3306

  kore-apiserver:
    image: quay.io/appvia/kore-apiserver:${KORE_TAG}
    environment:
      KORE_ADMIN_PASS: password
      KORE_ADMIN_TOKEN: password
      KORE_API_PUBLIC_URL: http://localhost:10080
      KORE_UI_PUBLIC_URL: http://localhost:3000
      KORE_AUTHENTICATION_PLUGINS: basicauth,admintoken,openid
      KORE_CERTIFICATE_AUTHORITY: ca/ca.pem
      KORE_CERTIFICATE_AUTHORITY_KEY: ca/ca-key.pem
      KUBE_API_SERVER: http://kube-apiserver:8080
      USERS_DB_URL: root:pass@tcp(database:3306)/kore?parseTime=true
      VERBOSE: 'true'
      KORE_IDP_CLIENT_ID: ${KORE_IDP_CLIENT_ID}
      KORE_IDP_CLIENT_SECRET: ${KORE_IDP_CLIENT_SECRET}
      KORE_IDP_SERVER_URL: ${KORE_IDP_SERVER_URL}
      KORE_IDP_USER_CLAIMS: preferred_username,email,name,username
      KORE_IDP_CLIENT_SCOPES: email,profile,offline_access
    ports:
      - 10080:10080
    restart: always
    # Used to source in the test certificate authority
    volumes:
      - ${KORE_LOCAL_HOME}/ca:/ca

  kore-ui:
    image: quay.io/appvia/kore-ui:${KORE_TAG}
    environment:
      KORE_BASE_URL: http://localhost:3000
      KORE_API_URL: http://kore-apiserver:10080/api/v1alpha1
      KORE_API_TOKEN: password
      REDIS_URL: redis://redis:6379
      KORE_IDP_CLIENT_ID: ${KORE_IDP_CLIENT_ID}
      KORE_IDP_CLIENT_SECRET: ${KORE_IDP_CLIENT_SECRET}
      KORE_IDP_SERVER_URL: ${KORE_IDP_SERVER_URL}
      KORE_IDP_USER_CLAIMS: preferred_username,email,name,username
      KORE_IDP_CLIENT_SCOPES: email,profile,offline_access
    ports:
      - 3000:3000
    restart: always

  redis:
    image: redis:5
    ports:
      - 6379:6379
    restart: always

  dex-operator:
    image: quay.io/appvia/dex:v2.20.0-master_grpc-connectors
    working_dir: /
    command:
      - serve
      - ./dex/config.yaml
    restart: always
    ports:
      - 5556:5556
      - 5557:5557
    volumes:
      - ${KORE_LOCAL_HOME}/dex:/dex
`

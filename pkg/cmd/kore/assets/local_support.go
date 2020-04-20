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

// LocalSupport indexes all of the required files for a local setup to run.
var LocalSupport = map[string]string{
	"dex/config.yaml":      LocalDexConfig,
	"dex/kubeconfig.local": LocalKubeConfigFile,
	"ca/ca.pem":            LocalDummyCAPEM,
	"ca/ca-key.pem":        LocalDummyCAPEMKey,
}

// LocalDexConfig is used by the local start command.
const LocalDexConfig = `
# DEX config file
issuer: http://localhost:5556/
storage:
  type: kubernetes
  config:
    kubeConfigFile: ./dex/kubeconfig.local
  waitForResources: true
web:
  http: 0.0.0.0:5556
oauth2:
  skipApprovalScreen: true
grpc:
  # Cannot be the same address as an HTTP(S) service.
  addr: 0.0.0.0:5557
  # Server certs. If TLS credentials aren't provided dex will run in plaintext (HTTP) mode.
  #tlsCert: /etc/dex/grpc.crt
  #tlsKey: /etc/dex/grpc.key
  # Client auth CA.
  #tlsClientCA: /etc/dex/client.crt
  # enable reflection
  reflection: true
enablePasswordDB: true
logger:
  level: "debug"
`

// LocalKubeConfigFile is used by the local start command.
const LocalKubeConfigFile = `
apiVersion: v1
clusters:
- cluster:
    server: http://kube-apiserver:8080
  name: local
contexts:
- context:
    cluster: local
    user: local
  name: local
current-context: local
kind: Config
preferences: {}
users:
- name: local
`

// LocalDummyCAPEM is used by the local start command.
const LocalDummyCAPEM = `
-----BEGIN CERTIFICATE-----
MIID6TCCAtGgAwIBAgIUUyy/2cR/bI6TxoE6UCdBHei7V7AwDQYJKoZIhvcNAQEL
BQAwgYMxCzAJBgNVBAYTAkdCMQswCQYDVQQIDAJHQjEPMA0GA1UEBwwGTG9uZG9u
MRMwEQYDVQQKDApBcHB2aWEgTHRkMQswCQYDVQQLDAJJVDEVMBMGA1UEAwwMY2Eu
YXBwdmlhLmlvMR0wGwYJKoZIhvcNAQkBFg5pbmZvQGFwcHZpYS5pbzAeFw0yMDAx
MjExMjUzMzlaFw0yODA0MDgxMjUzMzlaMIGDMQswCQYDVQQGEwJHQjELMAkGA1UE
CAwCR0IxDzANBgNVBAcMBkxvbmRvbjETMBEGA1UECgwKQXBwdmlhIEx0ZDELMAkG
A1UECwwCSVQxFTATBgNVBAMMDGNhLmFwcHZpYS5pbzEdMBsGCSqGSIb3DQEJARYO
aW5mb0BhcHB2aWEuaW8wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDN
iQatCOO9HdKyqFt/kqP/8wM3ZvcWwoacq+/YrwTHM07PyyG0+Lc+LoxBVtIVJUbA
yuTAma2rtTEEY1CbMPCrGfLXiSlALh/2HASgd/eO/XH0WgJorlnbJ8+Fr1dGjU8i
7Xn3Wo6WFvs2s2D+Qu1/KpaCqbKZ4ltB7CLz3kp0cSBzVKVenrwQHEeeQKD7rjY1
zuAZeagsKANpz9ksg1xqGFBiODJclWVBUr9Pti57oaXmlCL3C4biWcdWheUc8c+1
Ml9LSwx9cL3Wq6V7zcNj6M5oIVWjtzaCiVnJMep/uU+TlEIqnoPkcyfXSe6rs2v+
FRixU4vBuZKO0aR0TCO5AgMBAAGjUzBRMB0GA1UdDgQWBBSAsZSsPIt1To3WWK2q
sf+rN14JLjAfBgNVHSMEGDAWgBSAsZSsPIt1To3WWK2qsf+rN14JLjAPBgNVHRMB
Af8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQBUrroaM9EjtSNLI9jFX66Hne8g
VB0D7+C7O69lfisfy6wtBaYJOTaWWdb1BTWI0SzWCmKhxjconvVn043ypnfDoqSS
f8c7yYxnN/yBvgiotZ/QkXkweZR1/sYTK+pcOw5c6zH9eK0KsVKWtyFDwhYS/hA2
wymDlSBXYHD7/VjcUH4UYIqILDGjFEm11z5Wxzrxvokdx9KXj47NdmkLONV7yvNR
dzlB9KWIEipea4R3z12Mjrio//i9vkNNtpTVUKyPe5uKZvjEIsALqAMoi6ITiYpI
rPcsVye6WJkym/wBRdjXu/mgXIQ7bT5GLsw4STIhNpJtYQ/IghFRAco9W+q4
-----END CERTIFICATE-----
`

// LocalDummyCAPEMKey is used by the local start command.
const LocalDummyCAPEMKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAzYkGrQjjvR3Ssqhbf5Kj//MDN2b3FsKGnKvv2K8ExzNOz8sh
tPi3Pi6MQVbSFSVGwMrkwJmtq7UxBGNQmzDwqxny14kpQC4f9hwEoHf3jv1x9FoC
aK5Z2yfPha9XRo1PIu1591qOlhb7NrNg/kLtfyqWgqmymeJbQewi895KdHEgc1Sl
Xp68EBxHnkCg+642Nc7gGXmoLCgDac/ZLINcahhQYjgyXJVlQVK/T7Yue6Gl5pQi
9wuG4lnHVoXlHPHPtTJfS0sMfXC91qule83DY+jOaCFVo7c2golZyTHqf7lPk5RC
Kp6D5HMn10nuq7Nr/hUYsVOLwbmSjtGkdEwjuQIDAQABAoIBADixTCscYZz/heeL
srlMnHnz8PYuK4eWnoTGlEDDfeDoURvV3vVJCVpYgo1fQlFc19hD3rcVbKcJMn0Q
W+KCrE+1t5smFT/DuUMsVUZh8OH7HJyW20U+mkBuCbrJM5ydS6/JqzPEQcI6ko5z
ChT4JwRFngBqiH4TxrI3TSjRLt5Q+xY2uQpx3upPbgyyu8/bRUXtUtZbj1tVhXfJ
XPSWKd2RrNDaaIkw+5WX6byAwgymwg+4NAqSUNa2wYD5T5cVa0v8nJ/idyfj/N8C
B5y2bJcxBrMsHByOozU5XSQtfmV6qmFiYF6rOqLifb2S9tAT0676wqmVm0XMLbvG
2A97BdUCgYEA9eCsFBTcn6/PL9F9aEVabfJGjcG1ae/IOEp21wujNCxt+NQncICS
5M90GEn0fjc4dZxz4zigkwJZ8zHwk+XsrskD92FjwgRBg7GpIOT49gG+cGkibEMY
0odSVkHevqTATSeOE2cNjHSJg5JXVftxo2ZfjNPFVHeYf9LcqwUOAvsCgYEA1f8u
jBaizvrw/qWfdMcNDqT0xlwJquKtE4ptHxqWDpyyTz+tduhSs8lo80qL/nKmwNvm
9FJ1Ykr7a9xDK/bW3w0Q3seLyr6Wsgl+w22T8ArdSXHZb5eOJoWQ6B5ALCOUHBio
c0Nd2lVX91cxKDdV96BMivfoEtiPgyZjPTxjFdsCgYATuQ70mWvNH2QmOM6vc4i6
cwm3y0cLFWHhKg/4VgWkZL/5isMTIi0mT4HHhP8otLNBs+gT3PH8eN7QRDxBENt4
dcVsrZI7+O1sa+7eJZ/W0/L7v2M0fflaweIX6za74ilOxxJ9efG7R4nUVQPOcNn/
unGFsWMN0H4aGsb6rPAfywKBgQC294Hy4P++/KvE7hMSI4a0eLGYT+UsKLdWt8pp
B7A5OhzyyT0lJ6pecdy794cOvTR6PQqQ51fZ/MZPCHqeQmShPWipMfACH0Z1Xsz1
huEwIfnl6+O/F9PAd/7Xl9XCZ4EhLKwKMRUzsjiOEAzFl9p26KXJRAE269Z4if/b
wZ/udQKBgC/K5w3zOzdi8XVF5gkn/rcsvp2yRJpBR1sVcCpLV8uqiVwMTM4R114p
QDj0xfAUTv2AFcnL+JszB7XWwMEAD00rc3BP1Zb/ORNgpTYlW2+iGtEc4JBCTK58
kqm2kB+ymWKkKGI9JAXTFEOuuCS2pqvW+rlzs/9weJdn6JpNOn58
-----END RSA PRIVATE KEY-----
`

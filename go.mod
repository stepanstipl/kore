module github.com/appvia/kore

go 1.14

require (
	cloud.google.com/go v0.61.0 // indirect
	github.com/Azure/azure-sdk-for-go v44.0.0+incompatible
	github.com/Azure/go-autorest/autorest v0.11.0
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.0
	github.com/Azure/go-autorest/autorest/to v0.4.0
	github.com/Azure/go-autorest/autorest/validation v0.3.0 // indirect
	github.com/Masterminds/goutils v1.1.0 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/RoaringBitmap/roaring v0.4.21 // indirect
	github.com/ahmetb/gen-crd-api-reference-docs v0.2.0
	github.com/apparentlymart/go-cidr v1.0.1
	github.com/armon/go-proxyproto v0.0.0-20200108142055-f0b8253b1507
	github.com/aws/aws-sdk-go v1.29.31
	github.com/banzaicloud/k8s-objectmatcher v1.3.3
	github.com/blevesearch/bleve v0.8.1
	github.com/blevesearch/blevex v0.0.0-20190916190636-152f0fe5c040 // indirect
	github.com/blevesearch/go-porterstemmer v1.0.2 // indirect
	github.com/blevesearch/segment v0.0.0-20160915185041-762005e7a34f // indirect
	github.com/client9/misspell v0.3.4
	github.com/containerd/containerd v1.3.3 // indirect
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/couchbase/vellum v0.0.0-20190829182332-ef2e028c01fd // indirect
	github.com/cznic/b v0.0.0-20181122101859-a26611c4d92d // indirect
	github.com/cznic/mathutil v0.0.0-20181122101859-297441e03548 // indirect
	github.com/cznic/strutil v0.0.0-20181122101858-275e90344537 // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/denisenkom/go-mssqldb v0.0.0-20200620013148-b91950f658ec // indirect
	github.com/dexidp/dex v0.0.0-00010101000000-000000000000
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.4.2-0.20200213202729-31a86c4ab209 // indirect
	github.com/emicklei/go-restful v2.11.1+incompatible
	github.com/emicklei/go-restful-openapi v1.2.0
	github.com/etcd-io/bbolt v1.3.3 // indirect
	github.com/evanphx/json-patch v4.5.0+incompatible
	github.com/facebookgo/ensure v0.0.0-20160127193407-b4ab57deab51 // indirect
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/facebookgo/subset v0.0.0-20150612182917-8dac2c3c4870 // indirect
	github.com/felixge/httpsnoop v1.0.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-bindata/go-bindata v3.1.2+incompatible
	github.com/go-logr/zapr v0.1.1 // indirect
	github.com/go-openapi/errors v0.19.4
	github.com/go-openapi/inflect v0.19.0
	github.com/go-openapi/runtime v0.19.12
	github.com/go-openapi/spec v0.19.7
	github.com/go-openapi/strfmt v0.19.5
	github.com/go-openapi/swag v0.19.8
	github.com/go-openapi/validate v0.19.7
	github.com/go-resty/resty/v2 v2.3.0
	github.com/go-swagger/go-swagger v0.23.0
	github.com/golang-migrate/migrate v3.5.4+incompatible
	github.com/golangci/golangci-lint v1.27.0
	github.com/google/go-cmp v0.5.1 // indirect
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.7.4 // indirect
	github.com/hashicorp/go-version v1.2.0
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/huandu/xstrings v1.2.1 // indirect
	github.com/idubinskiy/schematyper v0.0.0-20190118213059-f71b40dac30d
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/jinzhu/gorm v1.9.12
	github.com/jmhodges/levigo v1.0.0 // indirect
	github.com/jpillora/backoff v1.0.0
	github.com/julienschmidt/httprouter v1.2.0
	github.com/kubernetes-sigs/go-open-service-broker-client v0.0.0-20200323235047-56a01c84bf43
	github.com/lib/pq v1.3.0 // indirect
	github.com/manifoldco/promptui v0.7.0
	github.com/maxbrunsfeld/counterfeiter/v6 v6.2.2
	github.com/mikefarah/yq/v3 v3.0.0-20200415014842-6f0a329331f9
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/mitchellh/gox v1.0.1
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/onsi/ginkgo v1.12.0
	github.com/onsi/gomega v1.9.0
	github.com/open-policy-agent/opa v0.21.1
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.1.0
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/procfs v0.0.8 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20190728182440-6a916e37a237 // indirect
	github.com/romanyx/polluter v1.2.2
	github.com/rs/xid v1.2.1
	github.com/satori/go.uuid v1.2.0
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749 // indirect
	github.com/shurcooL/vfsgen v0.0.0-20181202132449-6a9ea43bcacd
	github.com/sirupsen/logrus v1.5.0
	github.com/skratchdot/open-golang v0.0.0-20200116055534-eef842397966
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/steveyen/gtreap v0.0.0-20150807155958-0abe01ef9be2 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/syndtr/goleveldb v1.0.0 // indirect
	github.com/tcnksm/ghr v0.13.0
	github.com/tecbot/gorocksdb v0.0.0-20191019123150-400c56251341 // indirect
	github.com/tidwall/gjson v1.6.0
	github.com/tidwall/sjson v1.1.1
	github.com/ulule/limiter/v3 v3.5.0
	github.com/urfave/cli/v2 v2.1.1
	github.com/xeipuuv/gojsonschema v1.2.0
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sys v0.0.0-20200724161237-0e2f3a69832c // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	golang.org/x/tools v0.0.0-20200725200936-102e7d357031 // indirect
	gonum.org/v1/gonum v0.7.0 // indirect
	google.golang.org/api v0.29.0
	google.golang.org/genproto v0.0.0-20200726014623-da3ae01ef02d // indirect
	google.golang.org/grpc v1.30.0
	gopkg.in/yaml.v2 v2.2.8
	gotest.tools/v3 v3.0.2 // indirect
	k8s.io/api v0.18.2
	k8s.io/apiextensions-apiserver v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v0.18.2
	k8s.io/code-generator v0.18.2
	k8s.io/kube-openapi v0.0.0-20200121204235-bf4fb3bd569c
	sigs.k8s.io/application v0.8.3
	sigs.k8s.io/aws-iam-authenticator v0.5.0
	sigs.k8s.io/controller-runtime v0.6.0
	sigs.k8s.io/controller-tools v0.2.5
	sigs.k8s.io/yaml v1.2.0
)

// TODO: use github hosted dex
replace github.com/dexidp/dex => github.com/appvia/dex v0.0.0-20191213161401-b147340b9bc0

// This is the kubernetes-1.14.1-tools tag (which is the same as upstream kubernetes-1.14.1)
replace k8s.io/code-generator => github.com/appvia/kubernetes-code-generator v0.0.0-20200311145355-28f8f0159a26

replace github.com/kubernetes-sigs/go-open-service-broker-client => github.com/appvia/go-open-service-broker-client v0.0.0-20200505172434-f28d621ea14e

replace github.com/ahmetb/gen-crd-api-reference-docs => github.com/appvia/gen-crd-api-reference-docs v0.2.1-0.20200604183043-37e61fdd102c

replace github.com/idubinskiy/schematyper => github.com/appvia/schematyper v0.0.0-20200710151743-82d7f6f07f29

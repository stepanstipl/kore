SHELL = /bin/sh -e
NAME=kore-apiserver
AUTHOR ?= appvia
AUTHOR_EMAIL=kore@appvia.io
BUILD_TIME=$(shell date '+%s')
CURRENT_TAG=$(shell git tag --points-at HEAD)
DEPS=$(shell go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
DOCKER_IMAGES ?= kore-apiserver auth-proxy
GIT_SHA=$(shell git --no-pager describe --always --dirty)
GIT_LAST_TAG_SHA=$(shell git rev-list --tags='v[0.9]*.[0-9]*.[0-9]*' --max-count=1)
GIT_LAST_TAG=$(shell git describe --tags $(GIT_LAST_TAG_SHA))
HUB_APIS_SHA=$(shell cd ../kore-apis && git log | head -n 1 | cut -d' ' -f2)
GOVERSION ?= 1.12.7
HARDWARE=$(shell uname -m)
PACKAGES=$(shell go list ./...)
REGISTRY=quay.io
ROOT_DIR=${PWD}
VETARGS ?= -asmdecl -atomic -bool -buildtags -copylocks -methods -nilfunc -printf -rangeloops -unsafeptr
APIS ?= $(shell find pkg/apis -name "v*" -type d | sed -e 's/pkg\/apis\///' | sort | tr '\n' ' ')
UNAME := $(shell uname)
KORE_DOCS_PATH ?= ${GOPATH}/src/github.com/appvia/kore-docs
ifeq ($(UNAME), Darwin)
API_HOST_IN_DOCKER = host.docker.internal
endif
ifeq ($(UNAME), Linux)
API_HOST_IN_DOCKER = 127.0.0.1
endif
ifeq ($(USE_GIT_VERSION),true)
	# in CI or if we're pushing images for test
	ifeq ($(CURRENT_TAG),)
	VERSION ?= $(GIT_LAST_TAG)-$(GIT_SHA)
	else
	VERSION ?= $(CURRENT_TAG)
	endif
else
	# case of local build - must find upstream images
	# - note the version reported by --version always includes the git SHA
	VERSION ?= $(GIT_LAST_TAG)
endif
LFLAGS ?= -X github.com/appvia/kore/pkg/version.Tag=${GIT_LAST_TAG} -X github.com/appvia/kore/pkg/version.GitSHA=${GIT_SHA} -X github.com/appvia/kore/pkg/version.Compiled=${BUILD_TIME} -X github.com/appvia/kore/pkg/version.Release=${VERSION}
CLI_PLATFORMS=darwin linux windows
CLI_ARCHITECTURES=386 amd64
export GOFLAGS = -mod=vendor

.PHONY: test authors changelog build docker release check vet glide-install demo golangci-lint apis swagger images

default: build

golang:
	@echo "--> Go Version"
	@go version
	@echo "GOFLAGS: $$GOFLAGS"

generate-assets:
	@echo "--> Generating assets"
	@go generate ./pkg/kore/assets
	@go generate -tags=dev ./pkg/kore/assets
	@go generate ./pkg/security

check-generate-assets: generate-assets
	@if [ $$(git status --porcelain pkg/apiclient | wc -l) -gt 0 ]; then \
		echo "There are local changes after running 'make generate-assets'. Did you forget to run it?"; \
		git status --porcelain pkg/apiclient; \
		exit 1; \
	fi

build: golang
	@echo "--> Compiling the project ($(VERSION))"
	@mkdir -p bin
	@for binary in kore kore-apiserver auth-proxy; do \
		echo "--> Building $${binary} binary" ; \
		CGO_ENABLED=0 go build -ldflags "${LFLAGS}" -tags=jsoniter -o bin/$${binary} cmd/$${binary}/*.go || exit 1; \
	done

kore: golang
	@echo "--> Compiling the kore binary"
	@mkdir -p bin
	go build -ldflags "${LFLAGS}" -tags=jsoniter -o bin/kore cmd/kore/*.go

auth-proxy: golang
	@echo "--> Compiling the auth-proxy binary"
	@mkdir -p bin
	CGO_ENABLED=0 go build -ldflags "${LFLAGS}" -o bin/auth-proxy cmd/auth-proxy/*.go

auth-proxy-image: golang
	@echo "--> Build the auth-proxy docker image"
	CGO_ENABLED=0 docker build -t ${REGISTRY}/${AUTHOR}/auth-proxy:${VERSION} -f images/Dockerfile.auth-proxy .

auth-proxy-image-release: auth-proxy-image
	@echo "--> Pushing auth image"
	docker push ${REGISTRY}/${AUTHOR}/auth-proxy:${VERSION}

kore-apiserver: golang
	@echo "--> Compiling the kore-apiserver binary"
	@mkdir -p bin
	CGO_ENABLED=0 go build -ldflags "${LFLAGS}" -o bin/kore-apiserver cmd/kore-apiserver/*.go

kore-apiserver-image: golang
	@echo "--> Compiling the kore-apiserver image"
	docker build -t ${REGISTRY}/${AUTHOR}/kore-apiserver:${VERSION} -f images/Dockerfile.kore-apiserver .

kore-apiserver-image-local:
	@echo "--> Building the kore-apiserver image local"
	docker build -t ${REGISTRY}/${AUTHOR}/kore-apiserver:${VERSION} -f images/Dockerfile.kore-apiserver.local .

docker-build:
	@echo "--> Running docker"
	docker run --rm \
		-v ${ROOT_DIR}:/go/src/github.com/${AUTHOR}/${NAME} \
		-v ${HOME}/.ssh:/root/.ssh \
		-w /go/src/github.com/${AUTHOR}/${NAME} \
		-e GOPRIVATE=github.com/appvia \
		-e GO111MODULE=on \
		golang:${GOVERSION} \
		make in-docker-build

verify-circleci:
	@echo "--> Verifying the circleci config"
	@docker run -ti --rm -v ${PWD}:/workspace \
		-w /workspace circleci/circleci-cli \
		circleci config validate

images:
	@echo "--> Building docker images"
	@for name in ${DOCKER_IMAGES}; do \
		echo "--> Building docker image $${name}" ; \
		docker build -t ${REGISTRY}/${AUTHOR}/$${name}:${VERSION} -f images/Dockerfile.$${name} . ; \
	done

kind-image-dev:
	@echo "--> Building dev Docker image for Kind"
	docker build -t ${REGISTRY}/${AUTHOR}/kore-apiserver:dev -f images/Dockerfile.kore-apiserver.dev images
	kind load docker-image ${REGISTRY}/${AUTHOR}/kore-apiserver:dev --name kore

push-images:
	@echo "--> Pushing docker images"
	@for name in ${DOCKER_IMAGES}; do \
		echo "--> Pushing docker image $${name}" ; \
		docker push ${REGISTRY}/${AUTHOR}/$${name}:${VERSION} ; \
	done

package:
	@rm -rf ./release
	@mkdir ./release
	@$(MAKE) package-cli
	@$(MAKE) package-helm
	cd ./release && sha256sum * > kore.sha256sums

package-cli:
	@echo "--> Compiling CLI static binaries"
	CGO_ENABLED=0 go run github.com/mitchellh/gox -parallel=4 -arch="${CLI_ARCHITECTURES}" -os="${CLI_PLATFORMS}" -ldflags "-w ${LFLAGS}" -output=./release/{{.Dir}}-cli-{{.OS}}-{{.Arch}} ./cmd/kore/

package-helm:
	@echo "--> Patching version tag into helm charts"
	go run github.com/mikefarah/yq/v3 w -i charts/kore/values.yaml api.image ${REGISTRY}/${AUTHOR}/kore-apiserver
	go run github.com/mikefarah/yq/v3 w -i charts/kore/values.yaml api.version ${VERSION}
	go run github.com/mikefarah/yq/v3 w -i charts/kore/values.yaml ui.image ${REGISTRY}/${AUTHOR}/kore-ui
	go run github.com/mikefarah/yq/v3 w -i charts/kore/values.yaml ui.version ${VERSION}
	@echo "--> Packaging helm chart release"
	@helm package ./charts/kore --app-version ${VERSION} --version ${VERSION} -d ./release/
	@mv ./release/kore-${VERSION}.tgz ./release/kore-helm-chart-${VERSION}.tgz

push-release-packages:
	@echo "--> Pushing compiled CLI binaries and helm chart to draft release (requires github token set in .gitconfig or GITHUB_TOKEN env variable)"
	go run github.com/tcnksm/ghr -replace -draft -n "Release ${VERSION}" "${VERSION}" ./release

in-docker-build:
	@echo "--> Building in Docker"
	@git config --global url.git@github.com:.insteadOf https://github.com/
	@$(MAKE) test
	@$(MAKE) build

swagger: compose
	@echo "--> Retrieving the swagger api"
	@cd bin && { ./kore-apiserver --kube-api-server http://127.0.0.1:8080 2>/dev/null >/dev/null & echo $$! > $@; }
	@$(MAKE) swagger-json
	@$(MAKE) compose-down
	@$(MAKE) swagger-validate

swagger-json: api-wait
	@curl --retry 20 --retry-delay 5 -sSL http://127.0.0.1:10080/swagger.json | jq . -M > swagger.json
	@echo "--> Copy swagger JSON for UI to use for its auto-gen"
	@cp swagger.json ui/kore-api-swagger.json

swagger-validate:
	@echo "--> Validating the swagger api"
	@go run github.com/go-swagger/go-swagger/cmd/swagger validate swagger.json --skip-warnings

swagger-apiclient:
	@$(MAKE) swagger-json
	@echo "--> Creating API client based on the swagger definition"
	@rm -r pkg/apiclient/* >/dev/null || true
	@go run github.com/go-swagger/go-swagger/cmd/swagger generate client -q -f swagger.json -c pkg/apiclient -m pkg/apiclient/models

check-swagger-apiclient: swagger-apiclient
	@if [ $$(git status --porcelain pkg/apiclient | wc -l) -gt 0 ]; then \
		echo "There are local changes after running 'make swagger-apiclient'. Did you forget to run it?"; \
		git status --porcelain pkg/apiclient; \
		exit 1; \
	fi

in-docker-swagger:
	@echo "--> Swagger in Docker"
	curl --retry 50 --retry-delay 3 --retry-connrefused -sSL http://${API_HOST}:10080/swagger.json | jq . > swagger.json
	@$(MAKE) swagger-validate

docker-swagger-validate:
	@echo "--> Running docker to run swagger"
	docker run --net host --rm \
		-v ${PWD}:${PWD} \
		-w ${PWD} \
		-e API_HOST=$(API_HOST_IN_DOCKER) \
		quay.io/appvia/kore-apiserver-build:v0.0.1 \
		make in-docker-swagger

compose-db:
	@docker-compose \
		--file hack/compose/database.yml \
		pull
	@docker-compose \
		--file hack/compose/database.yml \
		up -d

compose-up:
	@docker-compose \
		--file hack/compose/database.yml \
		--file hack/compose/kube.yml \
		pull
	@docker-compose \
		--file hack/compose/database.yml \
		--file hack/compose/kube.yml \
		up -d

compose: build compose-up
	@echo "--> Building a test environment"
	@echo "--> Open a browser: http://localhost:3000"
	@echo "--> Note: the UI is running on host-network (some machines might have an issue)"
	@echo "--> Remember to start the api: make run"

.PHONY: run
run:
	@echo "--> Starting dependencies..."
	@$(MAKE) compose-up
	@$(MAKE) kube-api-wait
	@$(MAKE) run-api

.PHONY: run-with-kind
run-with-kind:
	@docker-compose \
		--file hack/compose/database.yml \
		pull
	@docker-compose \
		--file hack/compose/database.yml \
		up -d
	@$(MAKE) run-api-with-kind

.PHONY: run-api
run-api: kore-apiserver
	@echo "--> Starting api..."
	@hack/bin/run-api.sh

.PHONY: run-api-with-kind
run-api-with-kind: kore-apiserver
	@echo "--> Starting api..."
	@hack/bin/run-api-with-kind.sh

kube-api-wait:
	@echo "--> Waiting for Kube API..."
	@hack/bin/http_test.sh http://127.0.0.1:8080

api-wait:
	@echo "--> Waiting for API..."
	@hack/bin/http_test.sh http://127.0.0.1:10080

ui-wait:
	@echo "--> Waiting for UI..."
	@hack/bin/http_test.sh http://127.0.0.1:3000 100

compose-down:
	@echo "--> Removing the test environment"
	@docker-compose \
		--file hack/compose/database.yml \
		--file hack/compose/kube.yml \
		--file hack/compose/operators.yml down

compose-logs:
	@echo "--> Following logs for the test environment"
	@docker-compose \
		--file hack/compose/database.yml \
		--file hack/compose/kube.yml \
		--file hack/compose/operators.yml logs -f

demo:
	@if ! ls demo.env >/dev/null 2>&1 ; then \
		echo "Demo details not set, please run:" ; \
		echo "    cp ./hack/compose/demo.env.tmpl ./demo.env" ; \
		echo "See docs for values then run:" ; \
		echo "    vi ./demo.env" ; \
		exit 1; \
	fi
	@echo "--> Building the demo environment"
	@echo "--> Open a browser: http://localhost:3000"
	@docker-compose \
		--file hack/compose/database.yml \
		--file hack/compose/kube.yml \
		--file hack/compose/demo.yml \
		pull
	@docker-compose \
		--file hack/compose/database.yml \
		--file hack/compose/kube.yml \
		--file hack/compose/demo.yml \
		up --force-recreate --renew-anon-volumes

docker-release: images push-images

docker: images

# provides a consistent build environment with swagger, jq and curl
docker-builder-release:
	@echo "--> Releasing a builder image"
	@$(MAKE) docker-builder-build
	@docker push ${REGISTRY}/${AUTHOR}/${NAME}-build:${VERSION}

docker-builder-build:
	@echo "--> Building the docker image"
	docker build --build-arg GOVERSION=${GOVERSION} -f ./hack/build/Dockerfile -t ${REGISTRY}/${AUTHOR}/${NAME}-build:${VERSION} .

release: build
	mkdir -p release
	gzip -c bin/${NAME} > release/${NAME}_${VERSION}_linux_${HARDWARE}.gz
	rm -f release/${NAME}

clean:
	@echo "--> Cleaning up the environment"
	rm -rf ./bin 2>/dev/null
	rm -rf ./release 2>/dev/null

authors:
	@echo "--> Updating the AUTHORS"
	git log --format='%aN <%aE>' | sort -u > AUTHORS

dep-install:
	@echo "--> Installing dependencies"
	@dep ensure -v

vet:
	@echo "--> Running go vet $(VETARGS) $(PACKAGES)"
	@go vet $(VETARGS) $(PACKAGES)

gofmt:
	@echo "--> Running gofmt check"
	@if gofmt -s -l $$(go list -f '{{.Dir}}' ./...) | grep -q \.go ; then \
		echo "You need to run the make format, we have file unformatted"; \
		gofmt -s -l $$(go list -f '{{.Dir}}' ./...); \
		exit 1; \
	fi

format:
	@echo "--> Running go fmt"
	@gofmt -s -w $$(go list -f '{{.Dir}}' ./...)

bench:
	@echo "--> Running go bench"
	@go test -bench=. -benchmem

verify-licences:
	@echo "--> Verifying the licence headers"
	@hack/verify-licence.sh

coverage:
	@echo "--> Running go coverage"
	@go test -coverprofile cover.out
	@go tool cover -html=cover.out -o cover.html

spelling:
	@echo "--> Checking the spelling"
	@find . -name "*.go" -type f -not -path "./vendor/*" -not -path "./ui/node_modules/*" | xargs go run github.com/client9/misspell/cmd/misspell -error -source=go *.go
	@find . -name "*.md" -type f -not -path "./vendor/*" -not -path "./ui/node_modules/*" | xargs go run github.com/client9/misspell/cmd/misspell -error -source=text *.md

golangci-lint:
	@echo "--> Checking against the golangci-lint"
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint run --timeout 5m -j 2 ./...

check:
	@echo "--> Running code checkers"
	@$(MAKE) golang
	@$(MAKE) gofmt
	@$(MAKE) golangci-lint
	@$(MAKE) spelling
	@$(MAKE) vet
	@$(MAKE) verify-licences
	@$(MAKE) check-generate-assets

test:
	@echo "--> Running the tests"
	@go test --cover -v $(PACKAGES)

run-api-test:
	(cd ${ROOT_DIR}/pkg/apiserver; go test -tags=integration -ginkgo.v -vet=off)

api-test:
	@$(MAKE) swagger-apiclient
	@$(MAKE) run-api-test

all: test
	@echo "--> Performing all tests"
	@$(MAKE) bench
	@$(MAKE) coverage

changelog: release
	git log $(shell git tag | tail -n1)..HEAD --no-merges --format=%B >> changelog

apis: golang
	@echo "--> Generating Clientsets & Deepcopies"
	@${MAKE} deepcopy-gen
	@${MAKE} register-gen
	@${MAKE} crd-gen
	@${MAKE} schema-gen

check-apis: apis
	@${MAKE} check-api-sync

check-api-sync:
	@if [ $$(git status --porcelain | wc -l) -gt 0 ]; then \
		echo "There are local changes after running 'make apis'. Did you forget to run it?"; \
		git status --porcelain; \
		exit 1; \
	fi

deepcopy-gen:
	@echo "--> Generating the deepcopies"
	@hack/update-codegen.sh

schema-gen:
	@echo "--> Generating the CRD definitions"
	@go run github.com/go-bindata/go-bindata/go-bindata \
		-nocompress \
		-pkg register \
	    -nometadata \
		-o pkg/register/assets.go \
		-prefix deploy deploy/crds
	@gofmt -s -w pkg/register/assets.go

openapi-gen:
	@echo "--> Generating OpenAPI files"
	@echo "--> packages $(APIS)"
	@$(foreach api,$(APIS), \
		go run k8s.io/kube-openapi/cmd/openapi-gen -h hack/boilerplate.go.txt \
			--output-file-base zz_generated_openapi \
			-i github.com/appvia/kore/pkg/apis/$(api) \
			-p github.com/appvia/kore/pkg/apis/$(api); )

register-gen:
	@echo "--> Generating Schema register.go"
	@echo "--> packages $(APIS)"
	@$(foreach api,$(APIS), \
		go run k8s.io/code-generator/cmd/register-gen -h hack/boilerplate.go.txt \
			--output-file-base zz_generated_register \
			-i github.com/appvia/kore/pkg/apis/$(api) \
			-p github.com/appvia/kore/pkg/apis/$(api); )

crd-gen:
	@echo "--> Generating CRD deployment files"
	@mkdir -p deploy
	@rm -f deploy/crds/* 2>/dev/null || true
	@go run sigs.k8s.io/controller-tools/cmd/controller-gen crd:trivialVersions=true,preserveUnknownFields=false paths=./pkg/apis/...  output:dir=deploy/crds

check-release-notes:
	@echo "--> Verifying Release Notes"
	@${MAKE} generate-release-notes
	@if [ $$(git status --porcelain | wc -l) -gt 0 ]; then \
		echo "There are local changes after running 'make generate-release-notes'. Did you forget to run it?"; \
		git status --porcelain; \
		exit 1; \
	fi

generate-release-notes:
	@go run ./hack/build/tools/awesome-release-logger/main.go -r -t ${VERSION} -notag -derivetag '-rc([0-9])*' -o CHANGELOG.md

.PHONY: kind-dev
kind-dev: kind-apiserver
	scripts/kind_dev.sh

kind-dev-down:
	kind delete cluster --name kore

kind-apiserver:
	@echo "--> Compiling the kore-apiserver binary for kind"
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -ldflags "${LFLAGS}" -o bin/kore-apiserver-linux-amd64 cmd/kore-apiserver/*.go

kind-apiserver-image:
	@echo "--> Compiling the kore-apiserver and loading into kind"
	@${MAKE} VERSION=${USER} docker
	kind load docker-image ${REGISTRY}/${AUTHOR}/kore-apiserver:${USER} --name kore
	kubectl --context kind-kore -n kore rollout restart deployment kore-apiserver

kind-apiserver-start:
	@kubectl --context=kind-kore -n kore patch deployment kore-apiserver --patch '{"spec":{"replicas":1}}'

kind-apiserver-stop:
	@kubectl --context=kind-kore -n kore patch deployment kore-apiserver --patch '{"spec":{"replicas":0}}'

kind-apiserver-reload: kind-apiserver-stop kind-apiserver kind-apiserver-start

kind-apiserver-logs:
	@while true; do kubectl --context=kind-kore -n kore logs -f -l name=kore-apiserver || true; sleep 1; done

kind-admintoken:
	@echo `kubectl --context kind-kore -n kore get secret kore-api -o json | jq -r ".data.KORE_ADMIN_TOKEN" | base64 --decode`

kind-api-test:
	@export KORE_ADMIN_TOKEN="$(shell kubectl --context kind-kore -n kore get secret kore-api -o json | jq -r ".data.KORE_ADMIN_TOKEN" | base64 --decode)" && ${MAKE} run-api-test

generate-crd-reference:
	@if [ ! -d "${KORE_DOCS_PATH}" ]; then \
  		echo "${KORE_DOCS_PATH} directory does not exist"; \
  		exit 1; \
  	fi
	echo "Generating CRD reference documentation"
	go run github.com/ahmetb/gen-crd-api-reference-docs \
		-api-dir=./pkg/apis \
		-config=./hack/crd-reference-doc-gen/config.json \
		-template-dir=./hack/crd-reference-doc-gen/template \
		-out-file="${KORE_DOCS_PATH}/content/_generated/kore_crd_reference.html"

.PHONY: generate-schema-structs
generate-schema-structs:
	go generate ./pkg/clusterproviders/...

.PHONY: check-generate-schema-structs
check-generate-schema-structs: generate-schema-structs
	@if [ $$(git status --porcelain | wc -l) -gt 0 ]; then \
		echo "There are local changes after running 'make generate-schema-structs'. Did you forget to run it?"; \
		git status --porcelain; \
		exit 1; \
	fi

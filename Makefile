SHELL = /bin/sh -e
NAME=kore-apiserver
AUTHOR ?= appvia
AUTHOR_EMAIL=gambol99@gmail.com
BUILD_TIME=$(shell date '+%s')
CURRENT_TAG=$(shell git tag --points-at HEAD)
DEPS=$(shell go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
GIT_SHA=$(shell git --no-pager describe --always --dirty)
GIT_LAST_TAG_SHA=$(shell git rev-list --tags --max-count=1)
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
LFLAGS ?= -X github.com/appvia/kore/pkg/version.GitSHA=${GIT_SHA} -X github.com/appvia/kore/pkg/version.Compiled=${BUILD_TIME} -X github.com/appvia/kore/pkg/version.Release=${VERSION}

export GOFLAGS = -mod=vendor

.PHONY: test authors changelog build docker release cover vet glide-install demo golangci-lint apis swagger images

default: build

golang:
	@echo "--> Go Version"
	@go version
	@echo "GOFLAGS: $$GOFLAGS"

generate-clusterappman-manifests:
	@echo "--> Generating static manifests"
	@cp ./hack/generate/manifests_vfsdata.go ./pkg/clusterappman/
	@go generate ./pkg/clusterappman >/dev/null

build: golang generate-clusterappman-manifests
	@echo "--> Compiling the project ($(VERSION))"
	@mkdir -p bin
	@for binary in kore-apiserver korectl auth-proxy kore-clusterappman; do \
		echo "--> Building $${binary} binary" ; \
		go build -ldflags "${LFLAGS}" -tags=jsoniter -o bin/$${binary} cmd/$${binary}/*.go || exit 1; \
	done

korectl: golang deps
	@echo "--> Compiling the korectl binary"
	@mkdir -p bin
	go build -ldflags "${LFLAGS}" -tags=jsoniter -o bin/korectl cmd/korectl/*.go

cobractl: golang deps
	@echo "--> Compiling the cobractl binary"
	@mkdir -p bin
	go build -ldflags "${LFLAGS}" -tags=jsoniter -o bin/cobractl cmd/cobractl/*.go

auth-proxy: golang deps
	@echo "--> Compiling the auth-proxy binary"
	@mkdir -p bin
	go build -ldflags "${LFLAGS}" -o bin/auth-proxy cmd/auth-proxy/*.go

kore-apiserver: golang deps generate-clusterappman-manifests
	@echo "--> Compiling the kore-apiserver binary"
	@mkdir -p bin
	go build -ldflags "${LFLAGS}" -o bin/kore-apiserver cmd/kore-apiserver/*.go

kore-clusterappman: golang generate-clusterappman-manifests deps
	@echo "--> Compiling the kore-clusterappman binary"
	@mkdir -p bin
	go build -ldflags "${LFLAGS}" -o bin/kore-clusterappman cmd/kore-clusterappman/*.go

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

images:
	@echo "--> Building docker images"
	@for name in kore-apiserver auth-proxy; do \
		echo "--> Building docker image $${name}" ; \
		docker build -t ${REGISTRY}/${AUTHOR}/$${name}:${VERSION} -f images/Dockerfile.$${name} . ; \
	done

push-images:
	@echo "--> Pushing docker images"
	@for name in kore-apiserver auth-proxy; do \
		echo "--> Pushing docker image $${name}" ; \
		docker push ${REGISTRY}/${AUTHOR}/$${name}:${VERSION} ; \
	done


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
#@kill `cat $@ 2>/dev/null` 2>/dev/null && rm $@ 2>/dev/null

swagger-json: api-wait
	@curl --retry 20 --retry-delay 5 --retry-connrefused -sSL http://127.0.0.1:10080/swagger.json | jq > swagger.json

swagger-validate:
	@echo "--> Validating the swagger api"
	@go run github.com/go-swagger/go-swagger/cmd/swagger validate swagger.json --skip-warnings

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

compose-up:
	@docker-compose \
		--file hack/compose/kube.yml pull
	@docker-compose \
		--file hack/compose/kube.yml up -d

compose: build compose-up
	@echo "--> Building a test environment"
	@echo "--> Open a browser: http://localhost:3000"
	@echo "--> Note: the UI is running on host-network (some machines might have an issue)"
	@echo "--> Remember to start the api: make run"

run:
	@echo "--> Starting dependancies..."
	@$(MAKE) compose-up
	@$(MAKE) kube-api-wait
	@$(MAKE) run-api-only

run-api-only: kore-apiserver
	@echo "--> Starting api..."
	@hack/bin/run-api-with-env.sh

kube-api-wait:
	@echo "--> Waiting for Kube API..."
	@hack/bin/http_test.sh http://127.0.0.1:8080

api-wait:
	@echo "--> Waiting for API..."
	@hack/bin/http_test.sh http://127.0.0.1:10080

compose-down:
	@echo "--> Removing the test environment"
	@docker-compose \
		--file hack/compose/kube.yml \
		--file hack/compose/operators.yml down

compose-logs:
	@echo "--> Following logs for the test environment"
	@docker-compose \
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
		--file hack/compose/kube.yml \
		--file hack/compose/demo.yml \
		--file hack/compose/operators.yml \
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

deps:
	@echo "--> Installing build dependencies"

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
	@echo "--> Verifiying the licence headers"
	@hack/verify-licence.sh

coverage:
	@echo "--> Running go coverage"
	@go test -coverprofile cover.out
	@go tool cover -html=cover.out -o cover.html

cover:
	@echo "--> Running go cover"
	@go test --cover $(PACKAGES)

spelling:
	@echo "--> Checking the spelling"
	@find . -name "*.go" -type f -not -path "./vendor/*" | xargs go run github.com/client9/misspell/cmd/misspell -error -source=go *.go
	@find . -name "*.md" -type f -not -path "./vendor/*" | xargs go run github.com/client9/misspell/cmd/misspell -error -source=text *.md

golangci-lint:
	@echo "--> Checking against the golangci-lint"
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint run ./...

test: generate-clusterappman-manifests
	@echo "--> Running the tests"
	@if [ ! -d "vendor" ]; then \
		make deps; \
  	fi
	@go test -v $(PACKAGES)
	@$(MAKE) golang
	@$(MAKE) gofmt
	@$(MAKE) golangci-lint
	@$(MAKE) spelling
	@$(MAKE) vet
	@$(MAKE) cover
	@$(MAKE) verify-licences

all: test
	@echo "--> Performing all tests"
	@$(MAKE) bench
	@$(MAKE) coverage

changelog: release
	git log $(shell git tag | tail -n1)..HEAD --no-merges --format=%B >> changelog

apis: golang
	@echo "--> Generating Clientsets & Deepcopies"
	@rm -rf pkg/client 2>/dev/null
	@${MAKE} deepcopy-gen
	@${MAKE} openapi-gen
	@${MAKE} register-gen
	@${MAKE} crd-gen
	@${MAKE} schema-gen

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

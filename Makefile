NAME=kore-apiserver
AUTHOR ?= appvia
AUTHOR_EMAIL=gambol99@gmail.com
BUILD_TIME=$(shell date '+%s')
CURRENT_TAG=$(shell git tag --points-at HEAD)
DEPS=$(shell go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
GIT_SHA=$(shell git --no-pager describe --always --dirty)
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
ifeq ($(CURRENT_TAG),)
VERSION ?= $(GIT_SHA)
else
VERSION ?= $(CURRENT_TAG)
endif
LFLAGS ?= -X github.com/appvia/kore/pkg/version.GitSHA=${GIT_SHA} -X github.com/appvia/kore/pkg/version.Compiled=${BUILD_TIME} -X github.com/appvia/kore/pkg/version.Release=${VERSION}

.PHONY: test authors changelog build docker static release cover vet glide-install demo golangci-lint apis go-swagger swagger

default: build

golang:
	@echo "--> Go Version"
	@go version

generate-clusterman-manifests:
	@echo "--> Generating static manifests"
	@cp ./hack/generate/manifests_vfsdata.go ./pkg/clusterman/
	@go generate ./pkg/clusterman >/dev/null

build: golang generate-clusterman-manifests
	@echo "--> Compiling the project ($(VERSION))"
	@mkdir -p bin
	@for binary in kore-apiserver korectl auth-proxy kore-clusterman; do \
		echo "--> Building $${binary} binary" ; \
		go build -ldflags "${LFLAGS}" -tags=jsoniter -o bin/$${binary} cmd/$${binary}/*.go ; \
	done

static: golang generate-clusterman-manifests deps
	@echo "--> Compiling the static binaries ($(VERSION))"
	@mkdir -p bin
	@for binary in kore-apiserver korectl auth-proxy kore-clusterman; do \
		echo "--> Building $${binary} binary" ; \
		CGO_ENABLED=0 GOOS=linux go build -ldflags "${LFLAGS}" -tags=jsoniter -o bin/$${binary} cmd/$${binary}/*.go ; \
	done

korectl: golang deps
	@echo "--> Compiling the korectl binary"
	@mkdir -p bin
	GOOS=linux go build -ldflags "${LFLAGS}" -tags=jsoniter -o bin/korectl cmd/korectl/*.go

auth-proxy: golang deps
	@echo "--> Compiling the auth-proxy binary"
	@mkdir -p bin
	GOOS=linux go build -ldflags "${LFLAGS}" -o bin/auth-proxy cmd/auth-proxy/*.go

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

images: static
	@$(MAKE) images-only

images-only:
	@echo "--> Building docker images"
	@for name in kore-apiserver auth-proxy; do \
		echo "--> Building docker image $${name}" ; \
		docker build -t ${REGISTRY}/${AUTHOR}/$${name}:${VERSION} -f images/Dockerfile.$${name} . ; \
	done

in-docker-build:
	@echo "--> Building in Docker"
	@git config --global url.git@github.com:.insteadOf https://github.com/
	@$(MAKE) test
	@$(MAKE) static

swagger: compose
	@echo "--> Retrieving the swagger api"
	@cd bin && { ./kore-apiserver --kube-api-server http://127.0.0.1:8080 2>/dev/null >/dev/null & echo $$! > $@; }
	@$(MAKE) swagger-json
	@$(MAKE) compose-down
	@$(MAKE) swagger-validate
#@kill `cat $@ 2>/dev/null` 2>/dev/null && rm $@ 2>/dev/null

swagger-json:
	@curl --retry 20 --retry-delay 5 --retry-connrefused -sSL http://127.0.0.1:10080/swagger.json | jq > swagger.json

swagger-validate: go-swagger
	@echo "--> Validating the swagger api"
	@swagger validate swagger.json --skip-warnings

go-swagger:
	@echo "--> Installing go-swagger tools"
	@swagger version >/dev/null 2>&1; if [ $$? -ne 0 ]; then \
		echo "--> Installing the go-swagger tools"; \
		GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger; \
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

compose-up:
	@docker-compose \
		--file hack/compose/kube.yml \
		--file hack/compose/operators.yml pull
	@docker-compose \
		--file hack/compose/kube.yml \
		--file hack/compose/operators.yml up -d

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

run-api-only:
	@echo "--> Starting api..."
	@hack/bin/run-api-with-env.sh

kube-api-wait:
	@echo "--> Waiting for Kube API..."
	@curl \
		--retry 50 \
		--retry-delay 3 \
		--retry-connrefused \
		-sSL http://127.0.0.1:8080 >/dev/null 2>&1

api-wait:
	@echo "--> Waiting for API..."
	@which curl
	@curl \
		--retry 50 \
		--retry-delay 3 \
		--retry-connrefused \
		-sSL http://127.0.0.1:10080 >/dev/null 2>&1

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
		up --force-recreate

docker-release:
	@echo "--> Building a release image"
	@$(MAKE) static
	@$(MAKE) docker
	@docker push ${REGISTRY}/${AUTHOR}/${NAME}:${VERSION}

docker: static
	@echo "--> Building the docker image"
	docker build -t ${REGISTRY}/${AUTHOR}/${NAME}:${VERSION} .

# provides a consistent build environment with swagger, jq and curl
docker-builder-release:
	@echo "--> Releasing a builder image"
	@$(MAKE) docker-builder-build
	@docker push ${REGISTRY}/${AUTHOR}/${NAME}-build:${VERSION}

docker-builder-build:
	@echo "--> Building the docker image"
	docker build --build-arg GOVERSION=${GOVERSION} -f ./hack/build/Dockerfile -t ${REGISTRY}/${AUTHOR}/${NAME}-build:${VERSION} .

release: static
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
	@go tool vet 2>/dev/null ; if [ $$? -eq 3 ]; then \
		GO111MODULE=off go get golang.org/x/tools/cmd/vet; \
	fi
	@go vet $(VETARGS) $(PACKAGES)

gofmt:
	@echo "--> Running gofmt check"
	@gofmt -s -l . \
	    | grep -q \.go ; if [ $$? -eq 0 ]; then \
            echo "You need to runn the make format, we have file unformatted"; \
            gofmt -s -l .; \
            exit 1; \
	    fi

format:
	@echo "--> Running go fmt"
	@gofmt -s -w .

bench:
	@echo "--> Running go bench"
	@go test -bench=. -benchmem

coverage:
	@echo "--> Running go coverage"
	@go test -coverprofile cover.out
	@go tool cover -html=cover.out -o cover.html

cover:
	@echo "--> Running go cover"
	@go test --cover $(PACKAGES)

spelling:
	@echo "--> Checking the spelling"
	@which misspell 2>/dev/null ; if [ $$? -eq 1 ]; then \
		GO111MODULE=off go get -u github.com/client9/misspell/cmd/misspell; \
	fi
	@misspell -error *.go
	@misspell -error *.md

golangci-lint:
	@echo "--> Checking against the golangci-lint"
	@which golangci-lint 2>/dev/null ; if [ $$? -eq 1 ]; then \
		GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint; \
	fi
	@golangci-lint run ./...

test: generate-clusterman-manifests
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
	@which go-bindata  2>/dev/null ; if [ $$? -eq 1 ]; then \
		go get -u github.com/go-bindata/go-bindata/...; \
	fi
	@go-bindata \
		-nocompress \
		-pkg register \
	    -nometadata \
		-o pkg/register/assets.go \
		-prefix deploy deploy/crds

openapi-gen:
	@echo "--> Generating OpenAPI files"
	@echo "--> packages $(APIS)"
	@which openapi-gen  2>/dev/null ; if [ $$? -eq 1 ]; then \
		go get -u k8s.io/kube-openapi/cmd/openapi-gen; \
	fi
	@$(foreach api,$(APIS), \
		openapi-gen -h hack/boilerplate.go.txt \
			--output-file-base zz_generated_openapi \
			-i github.com/appvia/kore/pkg/apis/$(api) \
			-p github.com/appvia/kore/pkg/apis/$(api); )

register-gen:
	@echo "--> Generating Schema register.go"
	@echo "--> packages $(APIS)"
	@$(foreach api,$(APIS), \
		register-gen -h hack/boilerplate.go.txt \
			--output-file-base zz_generated_register \
			-i github.com/appvia/kore/pkg/apis/$(api) \
			-p github.com/appvia/kore/pkg/apis/$(api); )

crd-gen:
	@echo "--> Generating CRD deployment files"
	@which controller-gen  2>/dev/null ; if [ $$? -eq 1 ]; then \
		GO111MODULE=off go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.4; \
	fi
	@mkdir -p deploy
	@rm -f deploy/crds/* 2>/dev/null || true
	@controller-gen crd:trivialVersions=true,preserveUnknownFields=false paths=./pkg/apis/...  output:dir=deploy/crds

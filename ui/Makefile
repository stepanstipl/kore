NAME=kore-ui
AUTHOR ?= appvia
REGISTRY=quay.io
VERSION ?= latest
OPENAPIGEN_VERSION=v4.3.1

default: build

deps: 
	@echo "--> Preparing"
	npm install

build:
	@echo "--> Building"
	npm run build

check:
	@echo "--> Checking"
	npm run lint

test:
	@echo "--> Testing"
	npm run lint
	npm test

test-e2e:
	@echo "--> Running E2E Tests (headless)"
	npm run test-e2e

test-e2e-debug:
	@echo "--> Running E2E Tests (non-headless)"
	@export SHOW_BROWSER=true && npm run test-e2e

test-e2e-nightly-gke:
	@echo "--> Running nightly E2E Tests for GKE (headless)"
	@export SHOW_BROWSER=false && npm run test-e2e-nightly-gke

test-e2e-nightly-gke-debug:
	@echo "--> Running nightly E2E Tests for GKE (non-headless)"
	@export SHOW_BROWSER=true && npm run test-e2e-nightly-gke

test-e2e-nightly-eks:
	@echo "--> Running nightly E2E Tests for EKS (headless)"
	@export SHOW_BROWSER=false && npm run test-e2e-nightly-eks

test-e2e-nightly-eks-debug:
	@echo "--> Running nightly E2E Tests for EKS (non-headless)"
	@export SHOW_BROWSER=true && npm run test-e2e-nightly-eks

kind-e2e:
	@export KORE_ADMIN_PASS="$(shell kubectl --context kind-kore -n kore get secret kore-api -o json | jq -r ".data.KORE_ADMIN_PASS" | base64 --decode)" && ${MAKE} test-e2e

kind-e2e-debug:
	@export KORE_ADMIN_PASS="$(shell kubectl --context kind-kore -n kore get secret kore-api -o json | jq -r ".data.KORE_ADMIN_PASS" | base64 --decode)" && ${MAKE} test-e2e-debug

kind-e2e-nightly-gke:
	@export KORE_ADMIN_PASS="$(shell kubectl --context kind-kore -n kore get secret kore-api -o json | jq -r ".data.KORE_ADMIN_PASS" | base64 --decode)" && ${MAKE} test-e2e-nightly-gke

kind-e2e-nightly-gke-debug:
	@export KORE_ADMIN_PASS="$(shell kubectl --context kind-kore -n kore get secret kore-api -o json | jq -r ".data.KORE_ADMIN_PASS" | base64 --decode)" && ${MAKE} test-e2e-nightly-gke-debug

kind-e2e-nightly-eks:
	@export KORE_ADMIN_PASS="$(shell kubectl --context kind-kore -n kore get secret kore-api -o json | jq -r ".data.KORE_ADMIN_PASS" | base64 --decode)" && ${MAKE} test-e2e-nightly-eks

kind-e2e-nightly-eks-debug:
	@export KORE_ADMIN_PASS="$(shell kubectl --context kind-kore -n kore get secret kore-api -o json | jq -r ".data.KORE_ADMIN_PASS" | base64 --decode)" && ${MAKE} test-e2e-nightly-eks-debug

update-swagger:
	@echo "--> Updating unit test / auto-gen swagger (requires API to be running locally)"
	@curl --retry 5 --retry-delay 5 --retry-connrefused -sSL http://127.0.0.1:10080/swagger.json | jq > ./kore-api-swagger.json
	@if [ $$(git status --porcelain ./kore-api-swagger.json  | wc -l) -gt 0 ]; then \
		echo "There are local changes after running 'make update-swagger'. BE SURE TO COMMIT THESE CHANGES."; \
		git status --porcelain ./kore-api-swagger.json; \
	fi

check-kore-autogen:
	@$(MAKE) generate-kore-models
	@$(MAKE) check-kore-models

generate-kore-models:
	@echo "--> Updating kore API models"
	@rm -r lib/kore-api/model/*.js
	@docker create -v /local --name src openapitools/openapi-generator-cli:${OPENAPIGEN_VERSION} /bin/true
	@docker cp $$(pwd)/kore-api-swagger.json src:/local/kore-api-swagger.json
	@docker run --rm --volumes-from src -e JAVA_OPTS="-Xmx1024M -DloggerPath=conf/log4j.properties -DmodelDocs=false -DmodelTests=false -Dmodels" openapitools/openapi-generator-cli:${OPENAPIGEN_VERSION} generate \
		-i /local/kore-api-swagger.json \
		-g javascript \
		-o /local/lib/ \
		-p sourceFolder=kore-api,emitModelMethods=true,usePromises=true
	@docker cp src:/local/lib/kore-api $$(pwd)/lib
	@docker rm src

check-kore-models:
	@if [ $$(git status --porcelain ./lib/kore-api/model | wc -l) -gt 0 ]; then \
		echo "There are local changes after running 'generate-kore-models'. Did you forget to run it?"; \
		git status --porcelain  ./lib/kore-api/model; \
		exit 1; \
	fi

docker-release:
	@echo "--> Building a release image"
	@$(MAKE) docker
	@docker push ${REGISTRY}/${AUTHOR}/${NAME}:${VERSION}

docker:
	@echo "--> Building the docker image"
	docker build -t ${REGISTRY}/${AUTHOR}/${NAME}:${VERSION} --build-arg version=${VERSION} .

compose:
	@echo "--> Pulling images"
	@docker-compose pull
	@echo "--> Starting dependencies"
	@docker-compose up -d

compose-down:
	@echo "--> Stopping dependencies"
	@docker-compose down

compose-logs:
	@echo "--> Stopping dependencies"
	@docker-compose logs -f

run:
	@$(MAKE) compose
	@./scripts/run-with-env.sh npm run dev

run-prod:
	@$(MAKE) build
	@$(MAKE) compose
	@./scripts/run-with-env.sh npm start

run-prod-no-animations:
	@$(MAKE) build
	@$(MAKE) compose
	@export KORE_UI_DISABLE_ANIMATIONS="true" && ./scripts/run-with-env.sh npm start

run-docker:
	@$(MAKE) docker
	@$(MAKE) compose
	docker run --env-file ../demo.env \
		-e KORE_BASE_URL=http://localhost:3000 \
		-e KORE_API_URL=http://host.docker.internal:10080/api/v1alpha1 \
		-e KORE_API_TOKEN=password \
		-e REDIS_URL=redis://host.docker.internal:6379 \
		-p 3000:3000 \
		"${REGISTRY}/${AUTHOR}/${NAME}:${VERSION}"

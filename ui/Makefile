NAME=kore-ui
AUTHOR ?= appvia
REGISTRY=quay.io
VERSION ?= latest

default: build

build:
	@echo "--> Building"
	npm install
	npm run build

test:
	@echo "--> Testing"
	npm install
	npm run lint
	npm test

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
	@docker run --rm -v ${PWD}/:/local openapitools/openapi-generator-cli generate \
		-i /local/kore-api-swagger.json \
		-g javascript \
		-o /local/lib/ \
		-p sourceFolder=kore-api,emitModelMethods=true,usePromises=true \
		-DmodelDocs=false -DmodelTests=false -Dmodels

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
	docker build -t ${REGISTRY}/${AUTHOR}/${NAME}:${VERSION} .

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
	@./scripts/run-dev.sh

run-prod:
	@$(MAKE) build
	@$(MAKE) compose
	@./scripts/run-prod.sh

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

update-test-swagger:
	@echo "--> Updating unit test swagger (requires API to be running locally)"
	@curl --retry 5 --retry-delay 5 --retry-connrefused -sSL http://127.0.0.1:10080/swagger.json | jq > ./__tests__/kore-api-swagger.json

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

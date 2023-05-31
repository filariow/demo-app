REPOSITORY_REF ?= redhat.io/demo-soa
FRONTEND_IMAGE_REF ?= $(REPOSITORY_REF)/frontend:latest
CATALOG_IMAGE_REF ?= $(REPOSITORY_REF)/catalog:latest
CATALOG_INIT_IMAGE_REF ?= $(REPOSITORY_REF)/catalog-init:latest
ORDERS_IMAGE_REF ?= $(REPOSITORY_REF)/orders:latest
ORDERS_INIT_IMAGE_REF ?= $(REPOSITORY_REF)/orders-init:latest
ORDERS_EVENTS_CONSUMER_IMAGE_REF ?= $(REPOSITORY_REF)/orders-events-consumer:latest

DOCKER_BUILD_ARGS ?=

MANIFESTS_FOLDER ?= config/manifests

# Local development

.PHONY: run-dev-locally
run-dev-locally: clean-local-dev
	docker compose -f docker-compose.yaml -f docker-compose.dev.yaml up --build

.PHONY: clean-local-dev
clean-local-dev:
	docker compose -f docker-compose.yaml -f docker-compose.dev.yaml down

.PHONY: run-locally
run-locally:
	docker compose up --build

.PHONY: run-backend-locally
run-backend-locally:
	docker compose up --build catalog catalog-init orders orders-init sns orders-db orders-events-consumer

# Kubernetes

.PHONY: manifests
manifests:
	kustomize build $(MANIFESTS_FOLDER)

.PHONY: install
install:
	kubectl apply -k $(MANIFESTS_FOLDER)

.PHONY: uninstall
uninstall:
	kubectl delete -k $(MANIFESTS_FOLDER)

.PHONY: deploy-cert-manager
deploy-cert-manager:
	kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.9.1/cert-manager.yaml
	kubectl rollout status -n cert-manager deploy/cert-manager-webhook -w --timeout=120s

# Docker

.PHONY: docker-build-frontend
docker-build-frontend:
	docker build $(DOCKER_BUILD_ARGS) -t $(FRONTEND_IMAGE_REF) frontend/eshop

.PHONY: docker-build-catalog
docker-build-catalog:
	docker build $(DOCKER_BUILD_ARGS) -t $(CATALOG_IMAGE_REF) -f services/catalog/deploy/catalog/Dockerfile services/catalog

.PHONY: docker-build-catalog-init
docker-build-catalog-init:
	docker build $(DOCKER_BUILD_ARGS) -t $(CATALOG_INIT_IMAGE_REF) -f services/catalog/deploy/init/Dockerfile services/catalog

.PHONY: docker-build-orders
docker-build-orders:
	docker build $(DOCKER_BUILD_ARGS) -t $(ORDERS_IMAGE_REF) -f services/orders/deploy/orders/Dockerfile services/orders

.PHONY: docker-build-orders-events-consumer
docker-build-orders-events-consumer:
	docker build $(DOCKER_BUILD_ARGS) -t $(ORDERS_EVENTS_CONSUMER_IMAGE_REF) -f services/orders-events-consumer/Dockerfile services/orders-events-consumer

.PHONY: docker-build-orders-init
docker-build-orders-init:
	docker build $(DOCKER_BUILD_ARGS) -t $(ORDERS_INIT_IMAGE_REF) -f services/orders/deploy/init/Dockerfile services/orders

.PHONY: docker-build-all
docker-build-all: docker-build-frontend docker-build-catalog docker-build-catalog-init docker-build-orders docker-build-orders-init docker-build-orders-events-consumer

.PHONY: docker-push-frontend
docker-push-frontend: docker-build-frontend
	docker push $(FRONTEND_IMAGE_REF)

.PHONY: docker-push-catalog
docker-push-catalog: docker-build-catalog
	docker push $(CATALOG_IMAGE_REF)

.PHONY: docker-push-catalog-init
docker-push-catalog-init: docker-build-catalog-init
	docker push $(CATALOG_INIT_IMAGE_REF)

.PHONY: docker-push-orders
docker-push-orders: docker-build-orders
	docker push $(ORDERS_IMAGE_REF)

.PHONY: docker-push-orders-events-consumer
docker-push-orders-events-consumer: docker-build-orders-events-consumer
	docker push $(ORDERS_EVENTS_CONSUMER_IMAGE_REF)

.PHONY: docker-push-orders-init
docker-push-orders-init: docker-build-orders-init
	docker push $(ORDERS_INIT_IMAGE_REF)

.PHONY: docker-push-all
docker-push-all: docker-build-all docker-push-frontend docker-push-orders docker-push-orders-init docker-push-orders-events-consumer docker-push-catalog docker-push-catalog-init

# Code Generation

FRONTEND_CLIENTS_FOLDER=frontend/eshop/src/Clients/

.PHONY: generate-orders-client
generate-orders-client:
	LANG=typescript-axios; docker run --rm -v ${PWD}:/local:Z --user $(id -u):$(id -g) openapitools/openapi-generator-cli \
		generate -i /local/services/orders/apis/openapi.yaml -g $$LANG -o /local/$(FRONTEND_CLIENTS_FOLDER)/orders/

.PHONY: generate-catalog-client
generate-catalog-client:
	LANG=typescript-axios; docker run --rm -v ${PWD}:/local:Z --user $(id -u):$(id -g) openapitools/openapi-generator-cli \
		generate -i /local/services/catalog/apis/openapi.yaml -g $$LANG -o /local/$(FRONTEND_CLIENTS_FOLDER)/catalog/

.PHONY: generate-clients
generate-clients: generate-catalog-client generate-orders-client

# print variables
.PHONY: print-%
print-%:
	@echo $($*)


CATALOG_IMAGE_REF ?= sbo-demo/catalog:latest
CATALOG_INIT_IMAGE_REF ?= sbo-demo/catalog-init:latest

.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o bin/catalog cmd/api/*.go

.PHONY: build-image
build-image:
	docker build -t $(CATALOG_IMAGE_REF) -f deploy/catalog/Dockerfile .

.PHONY: build-init
build-init:
	CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o bin/catalog-init cmd/init/*.go

.PHONY: build-image-init
build-image-init:
	docker build -t $(CATALOG_INIT_IMAGE_REF) -f deploy/init/Dockerfile .


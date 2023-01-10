ORDERS_IMAGE_REF ?= sbo-demo/orders:latest
ORDERS_INIT_IMAGE_REF ?= sbo-demo/orders-init:latest

.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o bin/orders cmd/api/*.go

.PHONY: build-image
build-image:
	docker build -t $(ORDERS_IMAGE_REF) -f deploy/orders/Dockerfile .

.PHONY: build-init
build-init:
	CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o bin/orders-init cmd/init/*.go

.PHONY: build-image-init
build-image-init:
	docker build -t $(ORDERS_INIT_IMAGE_REF) -f deploy/init/Dockerfile .


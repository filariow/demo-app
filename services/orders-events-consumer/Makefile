IMAGE_TAG ?= sbo-demo/orders-events-consumer:latest

.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o bin/orders-events-consumer cmd/main.go

.PHONY: build-image
build-image:
	docker build -t $(IMAGE_TAG) .


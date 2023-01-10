FROM golang:1.19.3 as builder
WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY Makefile Makefile

COPY cmd/init/ cmd/init
COPY pkg/ pkg
RUN make build-init

FROM gcr.io/distroless/static:nonroot
WORKDIR /
USER 65532:65532
COPY --from=builder /workspace/bin/orders-init /orders-init

ENTRYPOINT ["/orders-init"]

FROM golang:1.19.3 as builder
WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY Makefile Makefile

COPY cmd/api cmd/api
COPY pkg/ pkg
RUN make build

FROM gcr.io/distroless/static:nonroot
WORKDIR /
USER 65532:65532
COPY --from=builder /workspace/bin/orders /orders

ENTRYPOINT ["/orders"]

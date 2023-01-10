FROM golang:1.19.3 as builder
WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY cmd/api/ cmd/api
COPY pkg/ pkg
COPY Makefile Makefile
RUN make build

FROM gcr.io/distroless/static:nonroot
WORKDIR /
USER 65532:65532
COPY --from=builder /workspace/bin/catalog /catalog

ENTRYPOINT ["/catalog"]

FROM golang:1.19.3 as builder
WORKDIR /workspace

COPY go.mod go.mod
RUN go mod download

COPY . .
RUN make build

FROM gcr.io/distroless/static:nonroot
WORKDIR /
USER 65532:65532
COPY --from=builder /workspace/bin/orders-events-consumer /orders-events-consumer

ENTRYPOINT ["/orders-events-consumer"]

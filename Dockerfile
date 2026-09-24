FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/hypercache ./cmd/hypercache

FROM alpine:3.20
RUN addgroup -S hypercache && adduser -S -G hypercache hypercache
COPY --from=builder /out/hypercache /usr/local/bin/hypercache
USER hypercache
EXPOSE 6379
ENTRYPOINT ["hypercache"]
CMD ["-addr=0.0.0.0:6379"]

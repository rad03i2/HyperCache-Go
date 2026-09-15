FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /hypercache cmd/hypercache/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /hypercache .
EXPOSE 6379 8080
CMD ["./hypercache", "-port=6379", "-http=8080"]

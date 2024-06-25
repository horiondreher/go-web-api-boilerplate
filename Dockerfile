# Build Stage
FROM golang:1.22.4-alpine3.19 AS builder

WORKDIR /app
COPY . .
RUN go build -o bin/api-boilerplate cmd/http/main.go

# Run Stage
FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/bin/api-boilerplate .

EXPOSE 8080

ENTRYPOINT [ "/app/api-boilerplate" ]
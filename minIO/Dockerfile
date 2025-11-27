FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app ./cmd/server

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /app /app
ENV MINIO_USE_SSL=false
EXPOSE 8080 50051
CMD ["/app"]

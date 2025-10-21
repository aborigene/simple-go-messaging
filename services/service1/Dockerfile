FROM golang:1.21-alpine AS builder

WORKDIR /build
COPY go.mod go.sum* ./
COPY main.go ./
RUN go mod tidy && go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o service1 main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /build/service1 ./
EXPOSE 8080
CMD ["./service1"]
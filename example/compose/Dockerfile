FROM golang:1.22-alpine AS builder
RUN apk add --no-cache --update go
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./cmd/app/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN adduser -D -g '' user
WORKDIR /app
COPY --from=builder /app/main .
RUN chown user:user /app/main
USER user
CMD ["./main"]

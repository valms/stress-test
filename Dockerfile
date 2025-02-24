# Primeira etapa: Build
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod .

COPY cmd/ cmd/
COPY service/ service/

RUN CGO_ENABLED=0 GOOS=linux go build -o stress-test ./cmd/main.go

# Segunda etapa: Imagem final
FROM alpine:latest

RUN apk add --no-cache ca-certificates &&  update-ca-certificates

RUN adduser -D appuser

WORKDIR /app

COPY --from=builder /app/stress-test .

RUN chown appuser:appuser /app/stress-test &&  chmod +x /app/stress-test

USER appuser

ENTRYPOINT ["./stress-test"]
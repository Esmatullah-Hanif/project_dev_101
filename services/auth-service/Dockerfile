FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY shared/ ./shared/
COPY services/auth-service/ ./services/auth-service/

RUN cd services/auth-service && go build -o /app/auth ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/auth .

EXPOSE 8001

CMD ["./auth"]

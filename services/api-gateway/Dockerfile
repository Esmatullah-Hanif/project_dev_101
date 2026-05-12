FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY shared/ ./shared/
COPY services/api-gateway/ ./services/api-gateway/

RUN cd services/api-gateway && go build -o /app/gateway ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/gateway .

EXPOSE 8000

CMD ["./gateway"]

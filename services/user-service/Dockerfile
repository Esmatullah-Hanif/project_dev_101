FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY shared/ ./shared/
COPY services/user-service/ ./services/user-service/

RUN cd services/user-service && go build -o /app/user ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/user .

EXPOSE 8002

CMD ["./user"]

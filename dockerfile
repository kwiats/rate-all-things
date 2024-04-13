FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
COPY .env .env

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o tit cmd/tit/main.go

# 2 Stage Run app
FROM alpine:latest  

WORKDIR /root/
COPY --from=builder /app/.env .
COPY --from=builder /app/tit .
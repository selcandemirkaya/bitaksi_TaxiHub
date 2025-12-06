# ---- Build stage ----
FROM golang:1.24-alpine AS builder

WORKDIR /app

# go.mod ve go.sum önce kopyalanır (cache için)
COPY go.mod go.sum ./
RUN go mod download

# tüm kaynak kodu kopyala
COPY . .

# binary build
RUN go build -o driver-service .

# ---- Run stage ----
FROM alpine:latest

WORKDIR /root/

# builder'dan binary al
COPY --from=builder /app/driver-service .

# servis portu
EXPOSE 8081

CMD ["./driver-service"]

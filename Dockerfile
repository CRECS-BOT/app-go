# ===== builder =====
FROM golang:1.25-alpine AS builder

WORKDIR /app

# consigliato per build più veloci
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# build del main in cmd/mybot
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o app-go ./cmd/crecs

# ===== runtime =====
FROM alpine:3.20

WORKDIR /app

# IMPORTANT: cert per HTTPS (Telegram, ecc.)
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/app-go ./app-go
COPY --from=builder /app/webapp ./webapp
COPY --from=builder /app/website ./website

EXPOSE 8080
CMD ["./app-go", "--port", "8080"]

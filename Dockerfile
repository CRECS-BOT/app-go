#app-go
FROM golang:1.22-alpine AS builder

WORKDIR /app

# copi tutto l'app-go (già con webapp/ e website/ pieni di statici)
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app-go .

# stage finale minimale
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/app-go .
COPY --from=builder /app/webapp ./webapp
COPY --from=builder /app/website ./website

EXPOSE 8080

CMD ["./app-go", "--port", "8080"]

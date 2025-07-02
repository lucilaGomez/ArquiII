# MODIFICAR hotel-service/Dockerfile

FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/server

FROM alpine:latest

# Instalar ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /app/

# Crear directorio para uploads
RUN mkdir -p uploads/hotels

# Copiar binario
COPY --from=builder /app/main .

# Exponer puerto
EXPOSE 8080

# Ejecutar aplicación
CMD ["./main"]
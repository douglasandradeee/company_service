FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# Instalar apenas o git e ca-certificates (necessários para builds)
RUN apk add --no-cache git ca-certificates

# Copiar mod files primeiro para melhor cache
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fonte
COPY . .

# Build da aplicação (binário estático)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-s -w' -o main ./cmd/api

# Imagem final mínima (sem vulnerabilidades)
FROM scratch

WORKDIR /app

# Copiar apenas o binário e certificados
COPY --from=builder /app/main .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Expor porta
EXPOSE 8080

# Comando de execução
CMD ["./main"]
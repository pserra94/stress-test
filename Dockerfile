# Estágio de build
FROM golang:1.21-alpine AS builder

# Instala certificados SSL necessários para requisições HTTPS
RUN apk --no-cache add ca-certificates git

# Define o diretório de trabalho
WORKDIR /app

# Copia os arquivos de dependências
COPY go.mod go.sum ./

# Baixa as dependências
RUN go mod download

# Copia o código fonte
COPY . .

# Compila a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o stresstest .

# Estágio final - imagem mínima
FROM alpine:latest

# Instala certificados SSL para requisições HTTPS
RUN apk --no-cache add ca-certificates

# Cria um usuário não-root
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Define o diretório de trabalho
WORKDIR /root/

# Copia o binário do estágio de build
COPY --from=builder /app/stresstest .

# Define o usuário para execução
USER appuser

# Define o comando padrão
ENTRYPOINT ["./stresstest"] 
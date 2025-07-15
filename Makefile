# Nome da aplicação
APP_NAME=stresstest
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Configurações de build
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Comandos padrão
.PHONY: help build test clean docker-build docker-run install deps fmt vet

# Ajuda - comando padrão
help: ## Exibe esta ajuda
	@echo "Comandos disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Build da aplicação
build: ## Compila a aplicação
	@echo "🔨 Compilando $(APP_NAME)..."
	@go build $(LDFLAGS) -o $(APP_NAME) .
	@echo "✅ Build concluído: $(APP_NAME)"

# Build para Windows
build-windows: ## Compila para Windows
	@echo "🔨 Compilando $(APP_NAME) para Windows..."
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(APP_NAME).exe .
	@echo "✅ Build Windows concluído: $(APP_NAME).exe"

# Build para Linux
build-linux: ## Compila para Linux
	@echo "🔨 Compilando $(APP_NAME) para Linux..."
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(APP_NAME)-linux .
	@echo "✅ Build Linux concluído: $(APP_NAME)-linux"

# Build para macOS
build-darwin: ## Compila para macOS
	@echo "🔨 Compilando $(APP_NAME) para macOS..."
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(APP_NAME)-darwin .
	@echo "✅ Build macOS concluído: $(APP_NAME)-darwin"

# Build para todas as plataformas
build-all: build-windows build-linux build-darwin ## Compila para todas as plataformas

# Instala dependências
deps: ## Baixa e instala dependências
	@echo "📦 Instalando dependências..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependências instaladas"

# Executa testes
test: ## Executa testes unitários
	@echo "🧪 Executando testes..."
	@go test -v ./...
	@echo "✅ Testes concluídos"

# Executa testes com cobertura
test-coverage: ## Executa testes com relatório de cobertura
	@echo "🧪 Executando testes com cobertura..."
	@go test -cover ./...
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Relatório de cobertura gerado: coverage.html"

# Executa benchmarks
benchmark: ## Executa testes de performance
	@echo "⚡ Executando benchmarks..."
	@go test -bench=. ./...

# Formata código
fmt: ## Formata o código Go
	@echo "🎨 Formatando código..."
	@go fmt ./...
	@echo "✅ Código formatado"

# Executa vet
vet: ## Executa go vet para análise estática
	@echo "🔍 Executando análise estática..."
	@go vet ./...
	@echo "✅ Análise concluída"

# Executa linter completo (se golangci-lint estiver instalado)
lint: ## Executa linter completo
	@echo "🔍 Executando linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint não encontrado, executando go vet..."; \
		$(MAKE) vet; \
	fi
	@echo "✅ Linter concluído"

# Build da imagem Docker
docker-build: ## Constrói a imagem Docker
	@echo "🐳 Construindo imagem Docker..."
	@docker build -t $(APP_NAME):$(VERSION) .
	@docker tag $(APP_NAME):$(VERSION) $(APP_NAME):latest
	@echo "✅ Imagem Docker construída: $(APP_NAME):$(VERSION)"

# Executa via Docker
docker-run: ## Executa um teste de exemplo via Docker
	@echo "🐳 Executando teste via Docker..."
	@docker run --rm $(APP_NAME):latest --url=https://httpbin.org/get --requests=50 --concurrency=5

# Executa bash no container
docker-shell: ## Abre shell no container Docker
	@docker run --rm -it --entrypoint=/bin/sh $(APP_NAME):latest

# Instala a aplicação no sistema
install: build ## Instala a aplicação no sistema
	@echo "📦 Instalando $(APP_NAME)..."
	@sudo cp $(APP_NAME) /usr/local/bin/
	@echo "✅ $(APP_NAME) instalado em /usr/local/bin/"

# Remove arquivos de build
clean: ## Remove arquivos gerados
	@echo "🧹 Limpando arquivos..."
	@rm -f $(APP_NAME) $(APP_NAME).exe $(APP_NAME)-* coverage.out coverage.html
	@docker image prune -f --filter label=stage=builder 2>/dev/null || true
	@echo "✅ Limpeza concluída"

# Executa teste de exemplo
example: build ## Executa um exemplo de teste
	@echo "🚀 Executando exemplo..."
	@./$(APP_NAME) --url=https://httpbin.org/get --requests=20 --concurrency=4

# Verifica qualidade do código
check: fmt vet test ## Executa formatação, vet e testes

# Pipeline completa de CI
ci: deps check build ## Pipeline completa de CI/CD

# Mostra informações da aplicação
info: ## Exibe informações sobre a aplicação
	@echo "📋 Informações da Aplicação:"
	@echo "  Nome: $(APP_NAME)"
	@echo "  Versão: $(VERSION)"
	@echo "  Build Time: $(BUILD_TIME)"
	@echo "  GOOS: $(GOOS)"
	@echo "  GOARCH: $(GOARCH)"
	@echo "  Go Version: $(shell go version)"

# Release - cria build para todas as plataformas
release: clean deps check build-all docker-build ## Prepara release completo
	@echo "🎉 Release $(VERSION) preparado!"
	@echo "Arquivos gerados:"
	@ls -la $(APP_NAME)*
	@echo "Imagem Docker: $(APP_NAME):$(VERSION)" 
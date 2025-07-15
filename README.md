# StressTest CLI

Uma ferramenta de linha de comando desenvolvida em Go para realizar testes de carga em serviços web.

## 🚀 Características

- **Testes de carga HTTP/HTTPS**: Suporte completo para requisições HTTP e HTTPS
- **Concorrência configurável**: Controle do número de requisições simultâneas
- **Relatórios detalhados**: Estatísticas completas com distribuição de códigos de status
- **Métricas de performance**: Tempo de resposta, throughput e estatísticas detalhadas
- **Containerizado**: Execução via Docker para facilidade de uso
- **Interrupção graceful**: Suporte a Ctrl+C para parar o teste a qualquer momento

## 📊 Métricas Reportadas

- Tempo total de execução
- Quantidade total de requisições realizadas
- Quantidade de requisições com status HTTP 200
- Distribuição completa de códigos de status HTTP
- Tempo médio, mínimo e máximo de resposta
- Requisições por segundo (throughput)
- Total de dados transferidos
- Resumo de erros (se houver)

## 🛠️ Instalação

### Opção 1: Executar com Docker (Recomendado)

```bash
# Build da imagem Docker
docker build -t stresstest .

# Executar teste
docker run stresstest --url=http://google.com --requests=1000 --concurrency=10
```

### Opção 2: Compilar localmente

```bash
# Clonar o repositório
git clone <repository-url>
cd stresstest

# Baixar dependências
go mod download

# Compilar
go build -o stresstest .

# Executar
./stresstest --url=http://google.com --requests=1000 --concurrency=10
```

## 📝 Uso

### Parâmetros Obrigatórios

- `--url`: URL do serviço a ser testado (deve incluir http:// ou https://)
- `--requests`: Número total de requisições (1 a 1.000.000)
- `--concurrency`: Número de requisições simultâneas (1 a 10.000)

### Exemplos de Uso

#### Teste básico com Docker
```bash
docker run stresstest --url=https://httpbin.org/get --requests=100 --concurrency=10
```

#### Teste de alta concorrência
```bash
docker run stresstest --url=https://jsonplaceholder.typicode.com/posts/1 --requests=5000 --concurrency=100
```

#### Teste local (se compilado localmente)
```bash
./stresstest --url=http://localhost:8080/api/health --requests=1000 --concurrency=50
```

## 📈 Exemplo de Saída

```
Iniciando teste de carga...
URL: https://httpbin.org/get
Requisições: 1000
Concorrência: 10
==================================================
Progresso: 10.0% (100/1000 requisições)
Progresso: 20.0% (200/1000 requisições)
...
Progresso: 100.0% (1000/1000 requisições)

============================================================
                 RELATÓRIO DE TESTE DE CARGA
============================================================

📊 RESUMO GERAL:
------------------------------
⏱️  Tempo total de execução: 5.234s
📈 Total de requisições: 1000
✅ Requisições bem-sucedidas: 995
❌ Requisições com falha: 5
📊 Taxa de sucesso: 99.50%
🚀 Requisições por segundo: 191.06 req/s
💾 Total de dados transferidos: 456.7 KB

🔍 DISTRIBUIÇÃO DE CÓDIGOS DE STATUS:
-----------------------------------
✅ HTTP 200 (Sucesso): 995 requisições (99.5%)
🚫 HTTP 0 (Erro de Conexão): 5 requisições (0.5%)

⚡ MÉTRICAS DE PERFORMANCE:
------------------------------
📊 Tempo médio de resposta: 52ms
🏃 Tempo mínimo de resposta: 28ms
🐌 Tempo máximo de resposta: 234ms
📏 Variação de tempo: 206ms

🚨 RESUMO DE ERROS:
-------------------------
❌ context deadline exceeded: 5 ocorrências

============================================================
Teste concluído com sucesso!
```

## 🏗️ Arquitetura

O projeto segue os padrões recomendados pela comunidade Go:

```
.
├── cmd/                    # Comandos CLI (Cobra)
│   └── root.go
├── internal/              # Código interno da aplicação
│   ├── models/           # Estruturas de dados
│   │   └── models.go
│   ├── stresstest/       # Lógica do teste de carga
│   │   └── executor.go
│   └── report/           # Formatação de relatórios
│       └── formatter.go
├── main.go               # Ponto de entrada
├── go.mod               # Dependências
├── Dockerfile           # Containerização
└── README.md            # Documentação
```

## 🧪 Executando Testes

```bash
# Executar testes unitários
go test ./...

# Executar testes com cobertura
go test -cover ./...

# Executar testes de performance
go test -bench=. ./...
```

## 🔧 Configurações Avançadas

### Variáveis de Ambiente

A aplicação respeita as seguintes variáveis de ambiente:

- `HTTP_TIMEOUT`: Timeout para requisições HTTP (padrão: 30s)
- `MAX_IDLE_CONNS`: Máximo de conexões idle (padrão: configuração do Go)

### Limites de Segurança

- Máximo de 1.000.000 requisições por teste
- Máximo de 10.000 requisições simultâneas
- Timeout de 30 segundos por requisição

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo `LICENSE` para mais detalhes.

## 🐛 Reportando Bugs

Para reportar bugs ou solicitar features, abra uma issue no repositório do projeto.

## ⚡ Performance

Esta ferramenta foi otimizada para:
- Baixo uso de memória através de pools de workers
- Alta concorrência com goroutines
- Coleta eficiente de métricas
- Relatórios formatados e informativos

## 🔒 Segurança

- Execução em container não-root
- Validação rigorosa de parâmetros de entrada
- Timeout configurável para evitar requisições infinitas
- Tratamento seguro de interrupções de sinal 
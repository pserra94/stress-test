#!/bin/bash

echo "🚀 Exemplos de uso do StressTest CLI"
echo "===================================="

# Verifica se o Docker está disponível
if command -v docker >/dev/null 2>&1; then
    echo "🐳 Docker encontrado - usando versão containerizada"
    STRESSTEST_CMD="docker run --rm stresstest:latest"
else
    echo "💻 Usando versão local"
    STRESSTEST_CMD="./stresstest"
fi

echo ""
echo "1️⃣ Teste básico com baixa concorrência:"
echo "Comando: $STRESSTEST_CMD --url=https://httpbin.org/get --requests=20 --concurrency=2"
$STRESSTEST_CMD --url=https://httpbin.org/get --requests=20 --concurrency=2

echo ""
echo "=================================="
echo ""

echo "2️⃣ Teste com média concorrência:"
echo "Comando: $STRESSTEST_CMD --url=https://jsonplaceholder.typicode.com/posts/1 --requests=50 --concurrency=10"
$STRESSTEST_CMD --url=https://jsonplaceholder.typicode.com/posts/1 --requests=50 --concurrency=10

echo ""
echo "=================================="
echo ""

echo "3️⃣ Teste de alta concorrência (comentado por segurança):"
echo "# $STRESSTEST_CMD --url=https://httpbin.org/delay/1 --requests=100 --concurrency=20"
echo "# Descomente a linha acima para executar um teste mais intensivo"

echo ""
echo "=================================="
echo ""

echo "4️⃣ Teste com site de exemplo:"
echo "Comando: $STRESSTEST_CMD --url=https://example.com --requests=30 --concurrency=5"
$STRESSTEST_CMD --url=https://example.com --requests=30 --concurrency=5

echo ""
echo "=================================="
echo ""

echo "✅ Exemplos concluídos!"
echo ""
echo "💡 Dicas:"
echo "- Ajuste os valores de --requests e --concurrency conforme necessário"
echo "- Use URLs de teste como httpbin.org para evitar sobrecarregar serviços reais"
echo "- Monitore o uso de recursos do sistema durante testes intensivos"
echo "- Use Ctrl+C para interromper um teste em andamento"
echo ""
echo "📖 Para mais informações: ./stresstest --help" 
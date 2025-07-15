package report

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"stresstest/internal/models"
)

// Formatter é responsável por formatar e exibir relatórios
type Formatter struct{}

// NewFormatter cria uma nova instância do formatador
func NewFormatter() *Formatter {
	return &Formatter{}
}

// PrintReport exibe o relatório completo do teste de carga
func (f *Formatter) PrintReport(result *models.StressTestResult) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("                 RELATÓRIO DE TESTE DE CARGA")
	fmt.Println(strings.Repeat("=", 60))

	f.printSummary(&result.Report)
	f.printStatusCodeDistribution(&result.Report)
	f.printPerformanceMetrics(&result.Report)
	f.printErrorSummary(result.Results)

	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("Teste concluído com sucesso!")
}

// printSummary exibe o resumo geral do teste
func (f *Formatter) printSummary(report *models.TestReport) {
	fmt.Println("\n📊 RESUMO GERAL:")
	fmt.Println(strings.Repeat("-", 30))

	fmt.Printf("⏱️  Tempo total de execução: %v\n", report.TotalTime.Round(time.Millisecond))
	fmt.Printf("📈 Total de requisições: %d\n", report.TotalRequests)
	fmt.Printf("✅ Requisições bem-sucedidas: %d\n", report.SuccessfulReqs)
	fmt.Printf("❌ Requisições com falha: %d\n", report.FailedReqs)

	if report.TotalRequests > 0 {
		successRate := float64(report.SuccessfulReqs) / float64(report.TotalRequests) * 100
		fmt.Printf("📊 Taxa de sucesso: %.2f%%\n", successRate)
	}

	fmt.Printf("🚀 Requisições por segundo: %.2f req/s\n", report.RequestsPerSec)
	fmt.Printf("💾 Total de dados transferidos: %s\n", f.formatBytes(report.TotalDataTransfer))
}

// printStatusCodeDistribution exibe a distribuição de códigos de status HTTP
func (f *Formatter) printStatusCodeDistribution(report *models.TestReport) {
	fmt.Println("\n🔍 DISTRIBUIÇÃO DE CÓDIGOS DE STATUS:")
	fmt.Println(strings.Repeat("-", 35))

	if len(report.StatusCodes) == 0 {
		fmt.Println("Nenhum código de status registrado")
		return
	}

	// Ordena os códigos de status para exibição consistente
	var codes []int
	for code := range report.StatusCodes {
		codes = append(codes, code)
	}
	sort.Ints(codes)

	for _, code := range codes {
		count := report.StatusCodes[code]
		percentage := float64(count) / float64(report.TotalRequests) * 100
		statusText := f.getStatusText(code)
		icon := f.getStatusIcon(code)

		fmt.Printf("%s HTTP %d (%s): %d requisições (%.1f%%)\n",
			icon, code, statusText, count, percentage)
	}
}

// printPerformanceMetrics exibe métricas de performance detalhadas
func (f *Formatter) printPerformanceMetrics(report *models.TestReport) {
	fmt.Println("\n⚡ MÉTRICAS DE PERFORMANCE:")
	fmt.Println(strings.Repeat("-", 30))

	fmt.Printf("📊 Tempo médio de resposta: %v\n", report.AvgResponseTime.Round(time.Millisecond))
	fmt.Printf("🏃 Tempo mínimo de resposta: %v\n", report.MinResponseTime.Round(time.Millisecond))
	fmt.Printf("🐌 Tempo máximo de resposta: %v\n", report.MaxResponseTime.Round(time.Millisecond))

	// Calcula estatísticas adicionais
	if report.MaxResponseTime > 0 {
		variation := report.MaxResponseTime - report.MinResponseTime
		fmt.Printf("📏 Variação de tempo: %v\n", variation.Round(time.Millisecond))
	}
}

// printErrorSummary exibe um resumo dos erros encontrados
func (f *Formatter) printErrorSummary(results []models.RequestResult) {
	errorCount := make(map[string]int)

	for _, result := range results {
		if result.Error != nil {
			errorCount[result.Error.Error()]++
		}
	}

	if len(errorCount) == 0 {
		return
	}

	fmt.Println("\n🚨 RESUMO DE ERROS:")
	fmt.Println(strings.Repeat("-", 25))

	for errorMsg, count := range errorCount {
		fmt.Printf("❌ %s: %d ocorrências\n", errorMsg, count)
	}
}

// getStatusText retorna uma descrição textual do código de status HTTP
func (f *Formatter) getStatusText(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "Sucesso"
	case code >= 300 && code < 400:
		return "Redirecionamento"
	case code >= 400 && code < 500:
		return "Erro do Cliente"
	case code >= 500:
		return "Erro do Servidor"
	case code == 0:
		return "Erro de Conexão"
	default:
		return "Desconhecido"
	}
}

// getStatusIcon retorna um ícone apropriado para o código de status
func (f *Formatter) getStatusIcon(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "✅"
	case code >= 300 && code < 400:
		return "🔄"
	case code >= 400 && code < 500:
		return "⚠️"
	case code >= 500:
		return "❌"
	case code == 0:
		return "🚫"
	default:
		return "❓"
	}
}

// formatBytes formata bytes em uma string legível (KB, MB, GB)
func (f *Formatter) formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := []string{"KB", "MB", "GB", "TB", "PB"}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}

// PrintQuickSummary exibe um resumo rápido para uso em logs
func (f *Formatter) PrintQuickSummary(report *models.TestReport) {
	fmt.Printf("Resumo: %d/%d requisições bem-sucedidas (%.1f%%) em %v (%.2f req/s)\n",
		report.SuccessfulReqs,
		report.TotalRequests,
		float64(report.SuccessfulReqs)/float64(report.TotalRequests)*100,
		report.TotalTime.Round(time.Millisecond),
		report.RequestsPerSec)
}

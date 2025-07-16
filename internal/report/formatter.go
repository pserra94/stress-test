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
	f.printErrorCluster(&result.Report)
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
	fmt.Println("\n🔍 DISTRIBUIÇÃO DETALHADA DE CÓDIGOS HTTP:")
	fmt.Println(strings.Repeat("-", 45))

	if len(report.StatusCodes) == 0 {
		fmt.Println("❌ Nenhum código de status registrado")
		return
	}

	// Agrupa códigos por categoria para melhor visualização
	categories := f.categorizeStatusCodes(report.StatusCodes)

	// Ordena os códigos dentro de cada categoria
	for category, codes := range categories {
		if len(codes) == 0 {
			continue
		}

		f.printStatusCategory(category, codes, report.TotalRequests)
	}

	// Exibe resumo por categoria
	f.printCategorySummary(report.StatusCodes, report.TotalRequests)
}

// categorizeStatusCodes organiza os códigos por categoria
func (f *Formatter) categorizeStatusCodes(statusCodes map[int]int) map[string][]StatusCodeInfo {
	categories := map[string][]StatusCodeInfo{
		"✅ SUCESSOS (2xx)":           {},
		"🔄 REDIRECIONAMENTOS (3xx)":  {},
		"⚠️  ERROS DO CLIENTE (4xx)": {},
		"❌ ERROS DO SERVIDOR (5xx)":  {},
		"🚫 ERROS DE CONEXÃO":         {},
		"❓ OUTROS":                   {},
	}

	for code, count := range statusCodes {
		info := StatusCodeInfo{
			Code:        code,
			Count:       count,
			Description: f.getDetailedStatusDescription(code),
		}

		switch {
		case code >= 200 && code < 300:
			categories["✅ SUCESSOS (2xx)"] = append(categories["✅ SUCESSOS (2xx)"], info)
		case code >= 300 && code < 400:
			categories["🔄 REDIRECIONAMENTOS (3xx)"] = append(categories["🔄 REDIRECIONAMENTOS (3xx)"], info)
		case code >= 400 && code < 500:
			categories["⚠️  ERROS DO CLIENTE (4xx)"] = append(categories["⚠️  ERROS DO CLIENTE (4xx)"], info)
		case code >= 500 && code < 600:
			categories["❌ ERROS DO SERVIDOR (5xx)"] = append(categories["❌ ERROS DO SERVIDOR (5xx)"], info)
		case code == 0:
			categories["🚫 ERROS DE CONEXÃO"] = append(categories["🚫 ERROS DE CONEXÃO"], info)
		default:
			categories["❓ OUTROS"] = append(categories["❓ OUTROS"], info)
		}
	}

	// Ordena códigos dentro de cada categoria
	for category := range categories {
		sort.Slice(categories[category], func(i, j int) bool {
			return categories[category][i].Code < categories[category][j].Code
		})
	}

	return categories
}

// StatusCodeInfo contém informações detalhadas sobre um código de status
type StatusCodeInfo struct {
	Code        int
	Count       int
	Description string
}

// printStatusCategory exibe uma categoria de códigos de status
func (f *Formatter) printStatusCategory(category string, codes []StatusCodeInfo, totalRequests int) {
	fmt.Printf("\n%s\n", category)

	for _, info := range codes {
		percentage := float64(info.Count) / float64(totalRequests) * 100

		// Formata a exibição com alinhamento melhor
		fmt.Printf("  📋 HTTP %d - %s\n", info.Code, info.Description)
		fmt.Printf("     📊 %d requisições (%.2f%%)\n", info.Count, percentage)

		// Adiciona barra de progresso visual para percentuais significativos
		if percentage > 1.0 {
			barLength := int(percentage / 5) // Cada caractere representa 5%
			if barLength > 20 {
				barLength = 20
			}
			bar := strings.Repeat("█", barLength)
			fmt.Printf("     📈 [%s%s]\n", bar, strings.Repeat("░", 20-barLength))
		}
		fmt.Println()
	}
}

// printCategorySummary exibe um resumo consolidado por categoria
func (f *Formatter) printCategorySummary(statusCodes map[int]int, totalRequests int) {
	fmt.Println("\n📈 RESUMO POR CATEGORIA:")
	fmt.Println(strings.Repeat("-", 30))

	summary := map[string]int{
		"✅ Sucessos (2xx)":          0,
		"🔄 Redirecionamentos (3xx)": 0,
		"⚠️  Erros Cliente (4xx)":   0,
		"❌ Erros Servidor (5xx)":    0,
		"🚫 Erros Conexão":           0,
		"❓ Outros":                  0,
	}

	for code, count := range statusCodes {
		switch {
		case code >= 200 && code < 300:
			summary["✅ Sucessos (2xx)"] += count
		case code >= 300 && code < 400:
			summary["🔄 Redirecionamentos (3xx)"] += count
		case code >= 400 && code < 500:
			summary["⚠️  Erros Cliente (4xx)"] += count
		case code >= 500 && code < 600:
			summary["❌ Erros Servidor (5xx)"] += count
		case code == 0:
			summary["🚫 Erros Conexão"] += count
		default:
			summary["❓ Outros"] += count
		}
	}

	for category, count := range summary {
		if count > 0 {
			percentage := float64(count) / float64(totalRequests) * 100
			fmt.Printf("%-25s %6d req (%.2f%%)\n", category, count, percentage)
		}
	}
}

// getDetailedStatusDescription retorna uma descrição detalhada para códigos HTTP específicos
func (f *Formatter) getDetailedStatusDescription(code int) string {
	descriptions := map[int]string{
		// 2xx Success
		200: "OK - Requisição bem-sucedida",
		201: "Created - Recurso criado com sucesso",
		202: "Accepted - Requisição aceita para processamento",
		204: "No Content - Sucesso sem conteúdo de resposta",

		// 3xx Redirection
		301: "Moved Permanently - Recurso movido permanentemente",
		302: "Found - Recurso encontrado em outro local",
		304: "Not Modified - Recurso não foi modificado",
		307: "Temporary Redirect - Redirecionamento temporário",
		308: "Permanent Redirect - Redirecionamento permanente",

		// 4xx Client Error
		400: "Bad Request - Requisição malformada",
		401: "Unauthorized - Autenticação necessária",
		403: "Forbidden - Acesso negado",
		404: "Not Found - Recurso não encontrado",
		405: "Method Not Allowed - Método HTTP não permitido",
		408: "Request Timeout - Timeout da requisição",
		409: "Conflict - Conflito na requisição",
		410: "Gone - Recurso não está mais disponível",
		429: "Too Many Requests - Muitas requisições (rate limit)",

		// 5xx Server Error
		500: "Internal Server Error - Erro interno do servidor",
		501: "Not Implemented - Funcionalidade não implementada",
		502: "Bad Gateway - Gateway inválido",
		503: "Service Unavailable - Serviço indisponível",
		504: "Gateway Timeout - Timeout do gateway",
		505: "HTTP Version Not Supported - Versão HTTP não suportada",

		// Connection Error
		0: "Erro de Conexão - Falha na conexão de rede",
	}

	if desc, exists := descriptions[code]; exists {
		return desc
	}

	// Fallback para códigos não mapeados
	switch {
	case code >= 200 && code < 300:
		return "Sucesso"
	case code >= 300 && code < 400:
		return "Redirecionamento"
	case code >= 400 && code < 500:
		return "Erro do Cliente"
	case code >= 500:
		return "Erro do Servidor"
	default:
		return "Código desconhecido"
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

	fmt.Println("\n🚨 RESUMO DETALHADO DE ERROS:")
	fmt.Println(strings.Repeat("-", 35))

	// Ordena erros por frequência (mais comuns primeiro)
	type errorInfo struct {
		message string
		count   int
	}

	var errors []errorInfo
	for errorMsg, count := range errorCount {
		errors = append(errors, errorInfo{message: errorMsg, count: count})
	}

	sort.Slice(errors, func(i, j int) bool {
		return errors[i].count > errors[j].count
	})

	for _, err := range errors {
		fmt.Printf("❌ %s\n", err.message)
		fmt.Printf("   📊 %d ocorrências\n", err.count)
		fmt.Println()
	}
}

// printErrorCluster exibe um cluster específico de erros agrupados
func (f *Formatter) printErrorCluster(report *models.TestReport) {
	// Coleta apenas códigos de erro (não-2xx)
	errorCodes := make(map[int]int)
	hasErrors := false

	for code, count := range report.StatusCodes {
		if code < 200 || code >= 300 {
			errorCodes[code] = count
			hasErrors = true
		}
	}

	if !hasErrors {
		return
	}

	fmt.Println("\n🚨 CLUSTER DE ERROS DETECTADOS:")
	fmt.Println(strings.Repeat("-", 40))

	// Ordena os códigos de erro para exibição consistente
	var codes []int
	for code := range errorCodes {
		codes = append(codes, code)
	}
	sort.Ints(codes)

	totalErrors := 0
	for _, count := range errorCodes {
		totalErrors += count
	}

	for _, code := range codes {
		count := errorCodes[code]
		percentage := float64(count) / float64(report.TotalRequests) * 100
		errorPercentage := float64(count) / float64(totalErrors) * 100

		icon := f.getStatusIcon(code)
		description := f.getShortErrorDescription(code)

		fmt.Printf("%s %d %s\n", icon, count, description)
		fmt.Printf("   📊 %.2f%% do total | %.1f%% dos erros\n", percentage, errorPercentage)

		// Adiciona barra visual para erros mais significativos
		if count > 1 {
			barLength := int(errorPercentage / 5) // Cada caractere representa 5% dos erros
			if barLength > 15 {
				barLength = 15
			}
			if barLength > 0 {
				bar := strings.Repeat("▓", barLength)
				fmt.Printf("   📈 [%s%s]\n", bar, strings.Repeat("░", 15-barLength))
			}
		}
		fmt.Println()
	}

	fmt.Printf("🔢 Total de erros: %d/%d requisições (%.2f%%)\n",
		totalErrors, report.TotalRequests,
		float64(totalErrors)/float64(report.TotalRequests)*100)
}

// getShortErrorDescription retorna uma descrição curta para o cluster de erros
func (f *Formatter) getShortErrorDescription(code int) string {
	switch code {
	// Erros de conexão
	case 0:
		return "Erros de Conexão"

	// 4xx Client Errors
	case 400:
		return "Erros Bad Request (400)"
	case 401:
		return "Erros de Autenticação (401)"
	case 403:
		return "Erros de Acesso Negado (403)"
	case 404:
		return "Erros Not Found (404)"
	case 405:
		return "Erros Method Not Allowed (405)"
	case 408:
		return "Erros de Timeout (408)"
	case 409:
		return "Erros de Conflito (409)"
	case 429:
		return "Erros Rate Limit (429)"

	// 5xx Server Errors
	case 500:
		return "Erros Internos do Servidor (500)"
	case 501:
		return "Erros Not Implemented (501)"
	case 502:
		return "Erros Bad Gateway (502)"
	case 503:
		return "Erros Service Unavailable (503)"
	case 504:
		return "Erros Gateway Timeout (504)"
	case 505:
		return "Erros HTTP Version (505)"

	// 3xx Redirects (pode ser considerado erro em alguns contextos)
	case 301:
		return "Redirect Permanente (301)"
	case 302:
		return "Redirect Temporário (302)"
	case 304:
		return "Not Modified (304)"
	case 307:
		return "Redirect Temporário (307)"
	case 308:
		return "Redirect Permanente (308)"

	default:
		if code >= 400 && code < 500 {
			return fmt.Sprintf("Erros do Cliente (%d)", code)
		} else if code >= 500 {
			return fmt.Sprintf("Erros do Servidor (%d)", code)
		} else if code >= 300 && code < 400 {
			return fmt.Sprintf("Redirecionamentos (%d)", code)
		} else {
			return fmt.Sprintf("Códigos Inesperados (%d)", code)
		}
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

package cmd

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"stresstest/internal/models"
	"stresstest/internal/report"
	"stresstest/internal/stresstest"

	"github.com/spf13/cobra"
)

var (
	targetURL   string
	requests    int
	concurrency int
)

// rootCmd representa o comando base quando chamado sem subcomandos
var rootCmd = &cobra.Command{
	Use:   "stresstest",
	Short: "Uma ferramenta CLI para testes de carga em serviços web",
	Long: `StressTest é uma ferramenta de linha de comando desenvolvida em Go 
para realizar testes de carga em serviços web.

A ferramenta permite especificar:
- URL do serviço a ser testado
- Número total de requisições
- Nível de concorrência (requisições simultâneas)

Exemplo de uso:
  stresstest --url=http://google.com --requests=1000 --concurrency=10`,
	RunE: runStressTest,
}

// Execute adiciona todos os comandos filhos ao comando root e configura as flags adequadamente
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Flags obrigatórias
	rootCmd.Flags().StringVar(&targetURL, "url", "", "URL do serviço a ser testado (obrigatório)")
	rootCmd.Flags().IntVar(&requests, "requests", 0, "Número total de requisições (obrigatório)")
	rootCmd.Flags().IntVar(&concurrency, "concurrency", 0, "Número de requisições simultâneas (obrigatório)")

	// Marca as flags como obrigatórias
	rootCmd.MarkFlagRequired("url")
	rootCmd.MarkFlagRequired("requests")
	rootCmd.MarkFlagRequired("concurrency")
}

// runStressTest executa o teste de carga principal
func runStressTest(cmd *cobra.Command, args []string) error {
	// Valida os parâmetros de entrada
	if err := validateParameters(); err != nil {
		return fmt.Errorf("parâmetros inválidos: %w", err)
	}

	// Cria o contexto com cancelamento para permitir interrupção graceful
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Configura o handler para capturar sinais de interrupção (Ctrl+C)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan
		fmt.Println("\n\n🛑 Interrupção detectada. Finalizando teste...")
		cancel()
	}()

	// Configura o teste
	config := models.TestConfig{
		URL:         targetURL,
		Requests:    requests,
		Concurrency: concurrency,
	}

	// Executa o teste
	executor := stresstest.NewExecutor()
	result, err := executor.Run(ctx, config)
	if err != nil {
		return fmt.Errorf("erro durante a execução do teste: %w", err)
	}

	// Exibe o relatório
	formatter := report.NewFormatter()
	formatter.PrintReport(result)

	return nil
}

// validateParameters valida os parâmetros de entrada
func validateParameters() error {
	// Valida URL
	if targetURL == "" {
		return fmt.Errorf("URL é obrigatória")
	}

	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return fmt.Errorf("URL inválida: %w", err)
	}

	if parsedURL.Scheme == "" {
		return fmt.Errorf("URL deve incluir o esquema (http:// ou https://)")
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("apenas esquemas HTTP e HTTPS são suportados")
	}

	// Valida número de requisições
	if requests <= 0 {
		return fmt.Errorf("número de requisições deve ser maior que 0")
	}

	if requests > 1000000 {
		return fmt.Errorf("número de requisições não pode exceder 1.000.000")
	}

	// Valida concorrência
	if concurrency <= 0 {
		return fmt.Errorf("nível de concorrência deve ser maior que 0")
	}

	if concurrency > 10000 {
		return fmt.Errorf("nível de concorrência não pode exceder 10.000")
	}

	if concurrency > requests {
		return fmt.Errorf("nível de concorrência não pode ser maior que o número total de requisições")
	}

	return nil
}

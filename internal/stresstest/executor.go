package stresstest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"stresstest/internal/models"
)

// Executor gerencia a execução do teste de carga
type Executor struct {
	client *http.Client
}

// NewExecutor cria uma nova instância do executor
func NewExecutor() *Executor {
	return &Executor{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Run executa o teste de carga com a configuração especificada
func (e *Executor) Run(ctx context.Context, config models.TestConfig) (*models.StressTestResult, error) {
	fmt.Printf("Iniciando teste de carga...\n")
	fmt.Printf("URL: %s\n", config.URL)
	fmt.Printf("Requisições: %d\n", config.Requests)
	fmt.Printf("Concorrência: %d\n", config.Concurrency)
	fmt.Println(strings.Repeat("=", 50))

	startTime := time.Now()

	// Canal para enviar trabalhos (requisições a serem feitas)
	jobs := make(chan int, config.Requests)

	// Canal para receber resultados
	results := make(chan models.RequestResult, config.Requests)

	// WaitGroup para esperar todos os workers terminarem
	var wg sync.WaitGroup

	// Inicia os workers
	for i := 0; i < config.Concurrency; i++ {
		wg.Add(1)
		go e.worker(ctx, config.URL, jobs, results, &wg)
	}

	// Envia todas as requisições para o canal de trabalhos
	go func() {
		defer close(jobs)
		for i := 0; i < config.Requests; i++ {
			select {
			case jobs <- i:
			case <-ctx.Done():
				return
			}
		}
	}()

	// Coleta todos os resultados
	var allResults []models.RequestResult
	var resultWg sync.WaitGroup
	resultWg.Add(1)

	go func() {
		defer resultWg.Done()
		for i := 0; i < config.Requests; i++ {
			select {
			case result := <-results:
				allResults = append(allResults, result)

				// Mostra progresso a cada 10% das requisições
				if (i+1)%(config.Requests/10) == 0 || i+1 == config.Requests {
					progress := float64(i+1) / float64(config.Requests) * 100
					fmt.Printf("Progresso: %.1f%% (%d/%d requisições)\n", progress, i+1, config.Requests)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// Espera todos os workers terminarem
	wg.Wait()
	close(results)

	// Espera todos os resultados serem coletados
	resultWg.Wait()

	totalTime := time.Since(startTime)

	// Gera o relatório
	report := e.generateReport(allResults, totalTime, config.Requests)

	return &models.StressTestResult{
		Config:  config,
		Report:  report,
		Results: allResults,
	}, nil
}

// worker executa requisições HTTP de forma concorrente
func (e *Executor) worker(ctx context.Context, url string, jobs <-chan int, results chan<- models.RequestResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case _, ok := <-jobs:
			if !ok {
				return // Canal fechado
			}
			result := e.makeRequest(ctx, url)
			results <- result
		case <-ctx.Done():
			return
		}
	}
}

// makeRequest executa uma única requisição HTTP
func (e *Executor) makeRequest(ctx context.Context, url string) models.RequestResult {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return models.RequestResult{
			StatusCode: 0,
			Duration:   time.Since(start),
			Error:      err,
		}
	}

	// Adiciona User-Agent para identificar o stress test
	req.Header.Set("User-Agent", "StressTest-CLI/1.0")

	resp, err := e.client.Do(req)
	duration := time.Since(start)

	if err != nil {
		return models.RequestResult{
			StatusCode: 0,
			Duration:   duration,
			Error:      err,
		}
	}
	defer resp.Body.Close()

	// Lê o corpo da resposta para calcular o tamanho
	bodyBytes, err := io.ReadAll(resp.Body)
	responseSize := int64(len(bodyBytes))

	if err != nil {
		return models.RequestResult{
			StatusCode:   resp.StatusCode,
			Duration:     duration,
			Error:        err,
			ResponseSize: responseSize,
		}
	}

	return models.RequestResult{
		StatusCode:   resp.StatusCode,
		Duration:     duration,
		Error:        nil,
		ResponseSize: responseSize,
	}
}

// generateReport cria o relatório consolidado dos resultados
func (e *Executor) generateReport(results []models.RequestResult, totalTime time.Duration, expectedRequests int) models.TestReport {
	report := models.TestReport{
		TotalTime:     totalTime,
		TotalRequests: len(results),
		StatusCodes:   make(map[int]int),
	}

	var totalDuration time.Duration
	var totalDataTransfer int64
	minDuration := time.Duration(^uint64(0) >> 1) // Max duration
	maxDuration := time.Duration(0)

	for _, result := range results {
		// Contabiliza códigos de status
		report.StatusCodes[result.StatusCode]++

		// Contabiliza sucessos (2xx) e falhas
		if result.StatusCode >= 200 && result.StatusCode < 300 && result.Error == nil {
			report.SuccessfulReqs++
		} else {
			report.FailedReqs++
		}

		// Calcula estatísticas de tempo de resposta
		totalDuration += result.Duration
		totalDataTransfer += result.ResponseSize

		if result.Duration < minDuration {
			minDuration = result.Duration
		}
		if result.Duration > maxDuration {
			maxDuration = result.Duration
		}
	}

	// Calcula médias e outras métricas
	if len(results) > 0 {
		report.AvgResponseTime = totalDuration / time.Duration(len(results))
		report.MinResponseTime = minDuration
		report.MaxResponseTime = maxDuration
		report.RequestsPerSec = float64(len(results)) / totalTime.Seconds()
		report.TotalDataTransfer = totalDataTransfer
	}

	return report
}

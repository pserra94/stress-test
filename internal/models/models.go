package models

import (
	"time"
)

// TestConfig contém a configuração para o teste de carga
type TestConfig struct {
	URL         string
	Requests    int
	Concurrency int
}

// RequestResult representa o resultado de uma requisição individual
type RequestResult struct {
	StatusCode   int
	Duration     time.Duration
	Error        error
	ResponseSize int64
}

// TestReport contém os resultados consolidados do teste
type TestReport struct {
	TotalTime         time.Duration
	TotalRequests     int
	SuccessfulReqs    int
	FailedReqs        int
	StatusCodes       map[int]int
	AvgResponseTime   time.Duration
	MinResponseTime   time.Duration
	MaxResponseTime   time.Duration
	RequestsPerSec    float64
	TotalDataTransfer int64
}

// StressTestResult encapsula todos os dados do teste
type StressTestResult struct {
	Config  TestConfig
	Report  TestReport
	Results []RequestResult
}

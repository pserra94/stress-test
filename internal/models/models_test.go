package models

import (
	"testing"
	"time"
)

func TestTestConfig(t *testing.T) {
	config := TestConfig{
		URL:         "https://example.com",
		Requests:    100,
		Concurrency: 10,
	}

	if config.URL != "https://example.com" {
		t.Errorf("Expected URL to be 'https://example.com', got '%s'", config.URL)
	}

	if config.Requests != 100 {
		t.Errorf("Expected Requests to be 100, got %d", config.Requests)
	}

	if config.Concurrency != 10 {
		t.Errorf("Expected Concurrency to be 10, got %d", config.Concurrency)
	}
}

func TestRequestResult(t *testing.T) {
	result := RequestResult{
		StatusCode:   200,
		Duration:     100 * time.Millisecond,
		Error:        nil,
		ResponseSize: 1024,
	}

	if result.StatusCode != 200 {
		t.Errorf("Expected StatusCode to be 200, got %d", result.StatusCode)
	}

	if result.Duration != 100*time.Millisecond {
		t.Errorf("Expected Duration to be 100ms, got %v", result.Duration)
	}

	if result.Error != nil {
		t.Errorf("Expected Error to be nil, got %v", result.Error)
	}

	if result.ResponseSize != 1024 {
		t.Errorf("Expected ResponseSize to be 1024, got %d", result.ResponseSize)
	}
}

func TestTestReport(t *testing.T) {
	statusCodes := make(map[int]int)
	statusCodes[200] = 95
	statusCodes[404] = 5

	report := TestReport{
		TotalTime:         5 * time.Second,
		TotalRequests:     100,
		SuccessfulReqs:    95,
		FailedReqs:        5,
		StatusCodes:       statusCodes,
		AvgResponseTime:   50 * time.Millisecond,
		MinResponseTime:   10 * time.Millisecond,
		MaxResponseTime:   200 * time.Millisecond,
		RequestsPerSec:    20.0,
		TotalDataTransfer: 102400,
	}

	if report.TotalRequests != 100 {
		t.Errorf("Expected TotalRequests to be 100, got %d", report.TotalRequests)
	}

	if report.SuccessfulReqs != 95 {
		t.Errorf("Expected SuccessfulReqs to be 95, got %d", report.SuccessfulReqs)
	}

	if report.FailedReqs != 5 {
		t.Errorf("Expected FailedReqs to be 5, got %d", report.FailedReqs)
	}

	if len(report.StatusCodes) != 2 {
		t.Errorf("Expected 2 status codes, got %d", len(report.StatusCodes))
	}

	if report.StatusCodes[200] != 95 {
		t.Errorf("Expected 95 requests with status 200, got %d", report.StatusCodes[200])
	}

	if report.StatusCodes[404] != 5 {
		t.Errorf("Expected 5 requests with status 404, got %d", report.StatusCodes[404])
	}
}

func TestStressTestResult(t *testing.T) {
	config := TestConfig{
		URL:         "https://example.com",
		Requests:    10,
		Concurrency: 2,
	}

	report := TestReport{
		TotalTime:     time.Second,
		TotalRequests: 10,
		StatusCodes:   make(map[int]int),
	}

	results := []RequestResult{
		{StatusCode: 200, Duration: 100 * time.Millisecond},
		{StatusCode: 200, Duration: 150 * time.Millisecond},
	}

	stressTestResult := StressTestResult{
		Config:  config,
		Report:  report,
		Results: results,
	}

	if stressTestResult.Config.URL != "https://example.com" {
		t.Errorf("Expected Config.URL to be 'https://example.com', got '%s'", stressTestResult.Config.URL)
	}

	if stressTestResult.Report.TotalRequests != 10 {
		t.Errorf("Expected Report.TotalRequests to be 10, got %d", stressTestResult.Report.TotalRequests)
	}

	if len(stressTestResult.Results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(stressTestResult.Results))
	}
}

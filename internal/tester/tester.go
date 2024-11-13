package tester

import (
	"burden/internal/config"
	"burden/internal/loader"
	"burden/internal/metrics"
	"burden/pkg/model"
	"log"
	"net/http"
	"time"
)

func RunTest(cfg *config.Config) *metrics.Metrics {
	var requests []model.Request
	var err error

	if cfg.CollectionFile != "" {
		requests, err = loader.LoadCollection(cfg.CollectionFile)
		if err != nil {
			log.Fatalf("Ошибка при загрузке коллекции: %v", err)
		}
	} else {
		requests = []model.Request{
			{Method: "GET", URL: cfg.URL},
		}
	}

	startTime := time.Now()
	completedRequests := 0
	var totalResponseTime, totalLatency float64
	errors := 0

	for i := 0; i < cfg.TotalRequests; i++ {
		start := time.Now()
		success := sendRequest(requests[i%len(requests)])
		elapsed := time.Since(start).Seconds()

		if success {
			completedRequests++
			totalResponseTime += elapsed
			totalLatency += elapsed / 2
		} else {
			errors++
		}
	}

	elapsedTime := time.Since(startTime).Seconds()
	throughput := float64(completedRequests) / elapsedTime
	avgResponseTime := totalResponseTime / float64(completedRequests)
	avgLatency := totalLatency / float64(completedRequests)

	errorRate := float64(errors) / float64(cfg.TotalRequests) * 100

	log.Printf("Test Completed. Throughput: %.2f RPS", throughput)

	return &metrics.Metrics{
		Throughput:    throughput,
		ResponseTime:  avgResponseTime,
		Latency:       avgLatency,
		Errors:        errors,
		TotalRequests: cfg.TotalRequests,
	}
}

func sendRequest(req model.Request) bool {
	client := &http.Client{}
	httpReq, err := http.NewRequest(req.Method, req.URL, nil)
	if err != nil {
		log.Printf("Make request failed: %v", err)
		return false
	}

	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	start := time.Now()
	resp, err := client.Do(httpReq)
	latency := time.Since(start).Seconds()

	if err != nil {
		log.Fatalf("Send request failed: %v", err)
		return false
	}
	defer resp.Body.Close()

	log.Printf("Successful request to %s. Latency: %.2f sec, Response code: %d", req.URL, latency, resp.StatusCode)
	return resp.StatusCode == http.StatusOK
}

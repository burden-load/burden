package tester

import (
	"burden/internal/config"
	"burden/internal/loader"
	"burden/internal/metrics"
	"burden/pkg/model"
	"log"
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
	totalRequests := cfg.TotalRequests
	completedRequests := 0

	for i := 0; i < totalRequests; i++ {
		// Затычка
		if sendRequest(requests[i%len(requests)]) {
			completedRequests++
		}
	}

	elapsedTime := time.Since(startTime).Seconds()
	throughput := float64(completedRequests) / elapsedTime

	log.Printf("Test Completed. Throughput: %.2f RPS", throughput)

	return &metrics.Metrics{
		Throughput: throughput,
	}
}

func sendRequest(req model.Request) bool {
	// Затычка
	time.Sleep(10 * time.Millisecond)
	return true
}

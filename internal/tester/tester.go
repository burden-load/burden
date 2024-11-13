package tester

import (
	"burden/internal/config"
	"burden/internal/loader"
	"burden/internal/metrics"
	"burden/pkg/model"
	"log"
	"net/http"
	"runtime"
	"sync"
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

	var wg sync.WaitGroup
	startTime := time.Now()
	completedRequests := 0
	errors := 0
	var totalResponseTime, totalLatency, downtime float64
	peakConcurrency := 0
	var memStats runtime.MemStats
	mu := sync.Mutex{}

	// Создаем канал для отслеживания завершения
	stopChannel := make(chan bool)

	if cfg.MaxErrors != nil {
		go func() {
			for {
				if errors > *cfg.MaxErrors {
					log.Printf("Превышен порог ошибок (%d). Завершение тестирования.", *cfg.MaxErrors)
					stopChannel <- true
					break
				}
				time.Sleep(1 * time.Second)
			}
		}()
	}

	for i := 0; i < cfg.Users; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			concurrency := 0
			for j := 0; j < cfg.TotalRequests/cfg.Users; j++ {
				select {
				case <-stopChannel:
					return
				default:
					mu.Lock()
					concurrency++
					if concurrency > peakConcurrency {
						peakConcurrency = concurrency
					}
					mu.Unlock()

					start := time.Now()
					success := sendRequest(requests[j%len(requests)])
					elapsed := time.Since(start).Seconds()

					if success {
						completedRequests++
						totalResponseTime += elapsed
						totalLatency += elapsed / 2
					} else {
						errors++
						downtime += elapsed
					}

					mu.Lock()
					concurrency--
					mu.Unlock()
				}
			}
		}()
	}

	wg.Wait()
	close(stopChannel)

	runtime.ReadMemStats(&memStats)
	memoryUsage := float64(memStats.Alloc) / 1024 / 1024

	elapsedTime := time.Since(startTime).Seconds()
	throughput := float64(completedRequests) / elapsedTime
	avgResponseTime := totalResponseTime / float64(completedRequests)
	avgLatency := totalLatency / float64(completedRequests)
	errorRate := float64(errors) / float64(cfg.TotalRequests) * 100
	resourceUtilization := memoryUsage / float64(memStats.Sys) * 100

	if cfg.Detailed {
		log.Printf("Детальный отчет: \nThroughput: %.2f req/sec\nСреднее время отклика: %.2f sec\nСредняя задержка: %.2f sec\nОшибки: %d (%.2f%%)\nПиковая нагрузка: %d\nDowntime: %.2f sec\nResource Utilization: %.2f%%", throughput, avgResponseTime, avgLatency, errors, errorRate, peakConcurrency, downtime, resourceUtilization)
	} else {
		log.Printf("Throughput: %.2f req/sec, Среднее время отклика: %.2f sec, Средняя задержка: %.2f sec", throughput, avgResponseTime, avgLatency)
	}

	return &metrics.Metrics{
		Throughput:          throughput,
		ResponseTime:        avgResponseTime,
		Latency:             avgLatency,
		Errors:              errors,
		TotalRequests:       cfg.TotalRequests,
		Concurrency:         cfg.Users,
		PeakLoad:            peakConcurrency,
		Downtime:            downtime,
		ResourceUtilization: resourceUtilization,
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

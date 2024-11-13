package tester

import (
	"burden/internal/config"
	"burden/internal/loader"
	"burden/internal/metrics"
	"burden/pkg/model"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"strings"
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
					success, response := sendRequest(requests[j%len(requests)])
					elapsed := time.Since(start).Seconds()

					if success {
						completedRequests++
						totalResponseTime += elapsed
						totalLatency += elapsed / 2
						log.Printf("Response data: %s", response)
					} else {
						errors++
						downtime += elapsed
						log.Printf("Request failed: %v", response)
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

func sendRequest(req model.Request) (bool, string) {
	client := &http.Client{
		Timeout: 30 * time.Second, // Устанавливаем таймаут на 30 секунд
	}

	// Формирование URL с параметрами
	urlWithParams := req.URL
	if len(req.Params) > 0 {
		urlWithParams += "?"
		for key, value := range req.Params {
			urlWithParams += fmt.Sprintf("%s=%s&", key, value)
		}
		urlWithParams = strings.TrimSuffix(urlWithParams, "&")
	}

	// Создание тела запроса, если оно есть
	var body io.Reader
	if req.Body != "" {
		body = strings.NewReader(req.Body)
	}

	// Создание нового HTTP-запроса
	httpReq, err := http.NewRequest(req.Method, urlWithParams, body)
	if err != nil {
		log.Printf("Error creating request: %v, Method: %s, URL: %s", err, req.Method, req.URL)
		return false, ""
	}

	// Установка заголовков, если они есть
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// Отправка запроса и замер времени
	var resp *http.Response
	var attempt int
	for attempt = 1; attempt <= 3; attempt++ { // Максимум 3 попытки
		start := time.Now()
		resp, err = client.Do(httpReq)
		latency := time.Since(start).Seconds()

		if err != nil {
			log.Printf("Attempt %d: Error sending request: %v, Method: %s, URL: %s", attempt, err, req.Method, req.URL)
			if attempt < 3 {
				log.Printf("Retrying request to %s", req.URL)
				time.Sleep(2 * time.Second) // Ожидаем 2 секунды перед повторной попыткой
				continue
			}
			log.Printf("Request failed after 3 attempts")
			return false, ""
		}

		// Успешный запрос
		defer resp.Body.Close()
		log.Printf("Request to %s successful. Latency: %.2f sec, Response code: %d", req.URL, latency, resp.StatusCode)
		break
	}

	// Чтение тела ответа
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return false, ""
	}

	// Проверяем успешные коды ответов (200-299)
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return true, string(responseBody)
	}

	// Логируем ошибку для других кодов ответа
	log.Printf("Request to %s failed with status code: %d, Response: %s", req.URL, resp.StatusCode, string(responseBody))
	return false, ""
}

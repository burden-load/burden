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
	"strings"
	"sync"
	"time"
)

func RunTest(cfg *config.Config) *metrics.Metrics {
	var requests []model.Request
	var err error

	// Загрузка запросов из файла коллекции или использования базового GET
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

	var (
		completedRequests int
		errors            int
		totalResponseTime float64
		totalLatency      float64
		downtime          float64
		startTime         = time.Now()
		sem               = make(chan struct{}, cfg.Users)
		mu                sync.Mutex
		peakConcurrency   int
	)

	// Определяем каналы для завершения теста
	stopChannel := make(chan struct{})

	// Если указано минимальное значение throughput
	if cfg.MinThroughput != nil && *cfg.MinThroughput > 0 {
		go func() {
			for {
				elapsedTime := time.Since(startTime).Seconds()
				if elapsedTime > 5 {
					currentThroughput := float64(completedRequests) / elapsedTime
					if currentThroughput < *cfg.MinThroughput {
						log.Fatalf("Throughput ниже минимального значения %.2f req/sec. Завершение тестирования.", *cfg.MinThroughput)
						close(stopChannel)
						return
					}
				}
				time.Sleep(1 * time.Second)
			}
		}()
	}

	// Если указан порог ошибок
	if cfg.MaxErrors != nil && *cfg.MaxErrors > 0 {
		go func() {
			for {
				if errors > *cfg.MaxErrors {
					log.Printf("Превышен порог ошибок (%d). Завершение тестирования.", *cfg.MaxErrors)
					close(stopChannel)
					return
				}
				time.Sleep(1 * time.Second)
			}
		}()
	}

	// Если указана продолжительность теста
	if cfg.TestDuration > 0 {
		go func() {
			time.Sleep(time.Duration(cfg.TestDuration) * time.Second)
			log.Println("Тест завершен по времени.")
			close(stopChannel)
		}()
	}

	var wg sync.WaitGroup
	requestsLimit := cfg.TotalRequests // Количество запросов по умолчанию
	if requestsLimit == 0 {
		requestsLimit = int(^uint(0) >> 1) // Устанавливаем максимальное значение int (если не указан totalRequests)
	}

	for i := 0; i < requestsLimit; i++ {
		select {
		case <-stopChannel:
			log.Println("Тест завершен по внешнему сигналу.")
			goto END
		default:
			sem <- struct{}{}
			wg.Add(1)
			go func(req model.Request) {
				defer func() {
					wg.Done()
					<-sem
				}()

				start := time.Now()
				success, response := sendRequest(req)
				elapsed := time.Since(start).Seconds()

				mu.Lock()
				if success {
					completedRequests++
					totalResponseTime += elapsed
					totalLatency += elapsed / 2
				} else {
					errors++
					downtime += elapsed
				}

				currentConcurrency := len(sem)
				if currentConcurrency > peakConcurrency {
					peakConcurrency = currentConcurrency
				}
				mu.Unlock()

				if success {
					log.Printf("Response: %s", response)
				}
			}(requests[i%len(requests)])
		}
	}

END:
	wg.Wait()

	elapsedTime := time.Since(startTime).Seconds()
	throughput := float64(completedRequests) / elapsedTime
	avgResponseTime := totalResponseTime / float64(completedRequests)
	avgLatency := totalLatency / float64(completedRequests)

	// Логирование результатов
	if cfg.Detailed {
		log.Printf("Детальный отчет:\nThroughput: %.5f req/sec\nСреднее время отклика: %.5f sec\nСредняя задержка: %.5f sec\nОшибки: %d\nПиковая конкуренция: %d\nDowntime: %.1f sec",
			throughput, avgResponseTime, avgLatency, errors, peakConcurrency, downtime)
	} else {
		log.Printf("Throughput: %.2f req/sec, Среднее время отклика: %.5f sec, Средняя задержка: %.5f sec",
			throughput, avgResponseTime, avgLatency)
	}

	return &metrics.Metrics{
		Throughput:    throughput,
		ResponseTime:  avgResponseTime,
		Latency:       avgLatency,
		Errors:        errors,
		TotalRequests: completedRequests,
		Concurrency:   cfg.Users,
		PeakLoad:      peakConcurrency,
		Downtime:      downtime,
	}
}

func sendRequest(req model.Request) (bool, string) {
	client := &http.Client{
		Timeout: 30 * time.Second,
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

	// Создание тела запроса
	var body io.Reader
	if req.Body != "" {
		body = strings.NewReader(req.Body)
	}

	// Создание HTTP-запроса
	httpReq, err := http.NewRequest(req.Method, urlWithParams, body)
	if err != nil {
		log.Printf("Ошибка создания запроса: %v, Method: %s, URL: %s", err, req.Method, req.URL)
		return false, ""
	}

	// Установка заголовков
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// Отправка запроса
	resp, err := client.Do(httpReq)
	if err != nil {
		log.Printf("Ошибка выполнения запроса: %v, Method: %s, URL: %s", err, req.Method, req.URL)
		return false, ""
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Ошибка чтения тела ответа: %v", err)
		return false, ""
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return true, string(responseBody)
	}

	log.Printf("Запрос к %s завершился с ошибкой: код ответа %d", req.URL, resp.StatusCode)
	return false, ""
}

func processRequests(
	requests []model.Request,
	cfg *config.Config,
	stopChannel <-chan bool,
	mu *sync.Mutex,
	completedRequests *int,
	totalResponseTime *float64,
	totalLatency *float64,
	peakConcurrency *int,
	errors *int,
	downtime *float64,
	wg *sync.WaitGroup) {

	defer wg.Done()

	concurrency := 0
	for j := 0; j < cfg.TotalRequests/cfg.Users; j++ {
		select {
		case <-stopChannel:
			return
		default:
			// Увеличиваем конкуренцию
			mu.Lock()
			concurrency++
			if concurrency > *peakConcurrency {
				*peakConcurrency = concurrency
			}
			mu.Unlock()

			// Начинаем измерение времени
			start := time.Now()
			success, response := sendRequest(requests[j%len(requests)])
			elapsed := time.Since(start).Seconds()

			// Обработка успешных и неуспешных запросов
			if success {
				*completedRequests++
				*totalResponseTime += elapsed
				*totalLatency += elapsed / 2
				log.Printf("Response data: %s", response)
			} else {
				*errors++
				*downtime += elapsed
				log.Printf("Request failed: %v", response)
			}

			// Уменьшаем конкуренцию
			mu.Lock()
			concurrency--
			mu.Unlock()
		}
	}
}

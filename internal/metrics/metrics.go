package metrics

import "log"

type Metrics struct {
	Throughput    float64 // Пропускная способность
	ResponseTime  float64 // Среднее время отклика
	Latency       float64 // Средняя задержка
	Errors        int     // Количество ошибок
	TotalRequests int     // Общее количество запросов
	Concurrency   int     // Одновременные запросы
	PeakLoad      int     // Пиковая нагрузка
	Downtime      float64 // Время простоя
}

func (metrics Metrics) Print(detailed bool) {
	log.Printf("Throughput: %.5f req/sec", metrics.Throughput)
	log.Printf("Среднее время отклика: %.5f sec", metrics.ResponseTime)
	log.Printf("Средняя задержка: %.5f sec", metrics.Latency)

	if detailed {
		log.Printf("Ошибки: %d", metrics.Errors)
		log.Printf("Конкурентные запросы: %d", metrics.Concurrency)
		log.Printf("Пиковая нагрузка: %d", metrics.PeakLoad)
		log.Printf("Время простоя: %.2f sec", metrics.Downtime)
	}
}

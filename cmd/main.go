package main

import (
	"burden/internal/config"
	"burden/internal/loader"
	"burden/internal/tester"
	"log"
)

func main() {
	cfg := config.ParseFlags()

	log.Println("Load testing starting")

	if cfg.CollectionFile != "" {
		err := loader.LoadCollection(cfg.CollectionFile)
		if err != nil {
			log.Fatalf("Ошибка при запуске коллекции: %v", err)
		}
	}

	metrics := tester.RunTest(cfg)

	log.Printf("Пропускная способность: %.2f", metrics.Throughput)
	log.Printf("Время отклика: %.2f", metrics.ResponseTime)
	log.Printf("Время задержки: %.2f", metrics.Latency)

	// Дополнительный вывод, если включен флаг verbose
	if cfg.Verbose {
		log.Printf("Ошибки: %d", metrics.Errors)
		log.Printf("Использование ресурсов: %.2f", metrics.ResourceUtilization)
		log.Printf("Конкурентные запросы: %d", metrics.Concurrency)
		log.Printf("Пиковая нагрузка: %.2f", metrics.PeakLoad)
		log.Printf("Время простоя: %.2f", metrics.Downtime)
	}
}

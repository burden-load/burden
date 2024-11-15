package config

import (
	"burden/internal/metrics"
	"fmt"
	"log"
	"os"
)

type Config struct {
	URL            string
	CollectionFile string
	Users          int
	TotalRequests  int
	TestDuration   int
	MaxErrors      *int
	Detailed       bool // Новый флаг для детализированного отчета
}

func SaveMetricsToGitHubOutput(metrics metrics.Metrics) {
	// Получаем путь к файлу GITHUB_OUTPUT из переменной окружения
	outputFile := os.Getenv("GITHUB_OUTPUT")
	if outputFile == "" {
		log.Fatalf("Переменная окружения GITHUB_OUTPUT не установлена")
	}

	// Открываем файл для записи
	file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Ошибка при открытии GITHUB_OUTPUT: %v", err)
	}
	defer file.Close()

	// Запись метрик в файл
	_, err = fmt.Fprintf(file, "throughput=%.2f\n", metrics.Throughput)
	if err != nil {
		log.Fatalf("Ошибка при записи в GITHUB_OUTPUT: %v", err)
	}

	_, err = fmt.Fprintf(file, "response_time=%.5f\n", metrics.ResponseTime)
	if err != nil {
		log.Fatalf("Ошибка при записи в GITHUB_OUTPUT: %v", err)
	}

	_, err = fmt.Fprintf(file, "latency=%.5f\n", metrics.Latency)
	if err != nil {
		log.Fatalf("Ошибка при записи в GITHUB_OUTPUT: %v", err)
	}

	_, err = fmt.Fprintf(file, "errors=%d\n", metrics.Errors)
	if err != nil {
		log.Fatalf("Ошибка при записи в GITHUB_OUTPUT: %v", err)
	}

	_, err = fmt.Fprintf(file, "concurrency=%d\n", metrics.Concurrency)
	if err != nil {
		log.Fatalf("Ошибка при записи в GITHUB_OUTPUT: %v", err)
	}

	_, err = fmt.Fprintf(file, "peak_load=%d\n", metrics.PeakLoad)
	if err != nil {
		log.Fatalf("Ошибка при записи в GITHUB_OUTPUT: %v", err)
	}

	_, err = fmt.Fprintf(file, "downtime=%.2f\n", metrics.Downtime)
	if err != nil {
		log.Fatalf("Ошибка при записи в GITHUB_OUTPUT: %v", err)
	}
}

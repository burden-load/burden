package config

import (
	"burden/internal/metrics"
	"errors"
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
	MinThroughput  *float64
}

// SaveMetricsToGitHubOutput сохраняет метрики в файл GITHUB_OUTPUT.
func SaveMetricsToGitHubOutput(metrics metrics.Metrics) error {
	// Получение пути к файлу GITHUB_OUTPUT
	outputFile := os.Getenv("GITHUB_OUTPUT")
	if outputFile == "" {
		return errors.New("переменная окружения GITHUB_OUTPUT не установлена")
	}

	// Открытие файла для записи
	file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("ошибка при открытии GITHUB_OUTPUT: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Printf("ошибка при закрытии файла: %v", cerr)
		}
	}()

	// Список метрик для записи
	metricsMap := map[string]interface{}{
		"throughput":    fmt.Sprintf("%.2f", metrics.Throughput),
		"response_time": fmt.Sprintf("%.5f", metrics.ResponseTime),
		"latency":       fmt.Sprintf("%.5f", metrics.Latency),
		"errors":        metrics.Errors,
		"concurrency":   metrics.Concurrency,
		"peak_load":     metrics.PeakLoad,
		"downtime":      fmt.Sprintf("%.2f", metrics.Downtime),
	}

	// Запись метрик в файл
	for key, value := range metricsMap {
		if err := writeToOutput(file, key, value); err != nil {
			return fmt.Errorf("ошибка при записи метрики '%s': %w", key, err)
		}
	}

	return nil
}

// writeToOutput записывает ключ-значение в указанный файл.
func writeToOutput(file *os.File, key string, value interface{}) error {
	_, err := fmt.Fprintf(file, "%s=%v\n", key, value)
	return err
}

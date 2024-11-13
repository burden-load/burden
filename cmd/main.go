package main

import (
	"burden/internal/config"
	"burden/internal/tester"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	// Определяем флаги командной строки
	url := flag.String("url", "", "URL для нагрузки (используется, если не указан файл коллекции)")
	collectionFile := flag.String("collection", "", "Путь к файлу коллекции запросов")
	users := flag.Int("users", 1, "Количество параллельных пользователей")
	totalRequests := flag.Int("requests", 100, "Общее количество запросов")
	maxErrors := flag.Int("max-errors", -1, "Максимально допустимое количество ошибок для остановки теста (-1 для отключения)")
	detailed := flag.Bool("detailed", false, "Выводить расширенные метрики")
	duration := flag.Int("duration", 10, "Длительность теста (по умолчанию 10s)")

	// Задаем пользовательскую функцию Usage для вывода справки
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n", os.Args[0])
		fmt.Println("Options:")

		// Пример ручного вывода флагов с двумя дефисами
		fmt.Println("  --url=<url>           URL для тестирования (обязательный)")
		fmt.Println("  --collection=<path>   Путь к коллекции (необязательный)")
		fmt.Println("  --users=<count>       Количество пользователей для нагрузки (по умолчанию 1)")
		fmt.Println("  --requests=<count>    Количество запросов (по умолчанию 1000)")
		fmt.Println("  --duration=<duration> Длительность теста (по умолчанию 10s)")
		fmt.Println("  --detailed            Выводить более подробные метрики")
		fmt.Println("  --max-errors=<count>  Максимальное количество ошибок до завершения с кодом 1")

		fmt.Println("\nПримеры использования:")
		fmt.Printf("  %s --url=http://example.com/api --users=10 --requests=1000\n", os.Args[0])
		fmt.Printf("  %s --collection=./example_collection.json --detailed\n", os.Args[0])
	}

	// Парсинг флагов
	flag.Parse()

	// Если не передано ни одного параметра, выводим справку и выходим
	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(1)
	}

	// Проверка обязательных параметров
	if *url == "" && *collectionFile == "" {
		log.Fatal("Необходимо указать либо URL, либо путь к файлу коллекции запросов")
	}

	// Создание конфигурации
	cfg := &config.Config{
		URL:            *url,
		CollectionFile: *collectionFile,
		Users:          *users,
		TotalRequests:  *totalRequests,
		Detailed:       *detailed,
		TestDuration:   *duration,
	}

	// Установка maxErrors только если он задан
	if *maxErrors >= 0 {
		cfg.MaxErrors = maxErrors
	}

	// Запуск теста
	log.Println("Запуск нагрузочного тестирования...")
	metrics := tester.RunTest(cfg)

	config.SaveMetricsToGitHubOutput(*metrics)

	// Вывод результатов
	log.Printf("Throughput: %.5f req/sec", metrics.Throughput)
	log.Printf("Среднее время отклика: %.5f sec", metrics.ResponseTime)
	log.Printf("Средняя задержка: %.5f sec", metrics.Latency)

	if cfg.Detailed {
		log.Printf("Ошибки: %d", metrics.Errors)
		log.Printf("Конкурентные запросы: %d", metrics.Concurrency)
		log.Printf("Пиковая нагрузка: %d", metrics.PeakLoad)
		log.Printf("Время простоя: %.2f sec", metrics.Downtime)
	}
}

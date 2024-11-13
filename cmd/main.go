package main

import (
	"burden/internal/config"
	"burden/internal/tester"
	"flag"
	"log"
)

func main() {
	// Определяем флаги командной строки
	url := flag.String("url", "", "URL для нагрузки (используется, если не указан файл коллекции)")
	collectionFile := flag.String("collection", "", "Путь к файлу коллекции запросов")
	users := flag.Int("users", 1, "Количество параллельных пользователей")
	totalRequests := flag.Int("requests", 100, "Общее количество запросов")
	maxErrors := flag.Int("max-errors", -1, "Максимально допустимое количество ошибок для остановки теста (-1 для отключения)")
	detailed := flag.Bool("detailed", false, "Выводить расширенные метрики")

	// Парсинг флагов
	flag.Parse()

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
	}

	// Установка maxErrors только если он задан
	if *maxErrors >= 0 {
		cfg.MaxErrors = maxErrors
	}

	// Запуск теста
	log.Println("Запуск нагрузочного тестирования...")
	metrics := tester.RunTest(cfg)

	// Вывод результатов
	log.Printf("Throughput: %.2f req/sec", metrics.Throughput)
	log.Printf("Среднее время отклика: %.2f sec", metrics.ResponseTime)
	log.Printf("Средняя задержка: %.2f sec", metrics.Latency)

	if cfg.Detailed {
		log.Printf("Ошибки: %d", metrics.Errors)
		log.Printf("Использование ресурсов: %.2f%%", metrics.ResourceUtilization)
		log.Printf("Конкурентные запросы: %d", metrics.Concurrency)
		log.Printf("Пиковая нагрузка: %d", metrics.PeakLoad)
		log.Printf("Время простоя: %.2f sec", metrics.Downtime)
	}
}

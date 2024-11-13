package config

import "flag"

type Config struct {
	URL            string
	CollectionFile string
	Users          int
	TotalRequests  int
	TestDuration   int
	MaxErrors      *int
	Detailed       bool // Новый флаг для детализированного отчета
}

func LoadConfig() *Config {
	url := flag.String("url", "http://localhost", "URL для нагрузки")
	collectionFile := flag.String("collection", "", "Файл с коллекцией запросов")
	users := flag.Int("users", 1, "Количество пользователей")
	totalRequests := flag.Int("requests", 100, "Общее количество запросов")
	testDuration := flag.Int("duration", 10, "Длительность теста в секундах")
	detailed := flag.Bool("detailed", false, "Вывод подробного отчета")

	flag.Parse()

	return &Config{
		URL:            *url,
		CollectionFile: *collectionFile,
		Users:          *users,
		TotalRequests:  *totalRequests,
		TestDuration:   *testDuration,
		Detailed:       *detailed,
	}
}

package config

import (
	"flag"
	"time"
)

type Config struct {
	URL            string
	CollectionFile string
	Users          int
	TotalRequests  int
	Duration       time.Duration
	Verbose        bool
}

func ParseFlags() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.URL, "url", "", "URL для нагрузки (необязательно, если задана коллекция)")
	flag.StringVar(&cfg.CollectionFile, "collection", "", "Путь к файлу коллекции запросов (заменяет URL)")
	flag.IntVar(&cfg.Users, "users", 1, "Количество параллельных пользователей")
	flag.IntVar(&cfg.TotalRequests, "requests", 0, "Общее число запросов (если задано, длительность игнорируется)")
	flag.DurationVar(&cfg.Duration, "duration", 0, "Длительность теста (если задана, количество запросов игнорируется)")
	flag.BoolVar(&cfg.Verbose, "verbose", false, "Вывод всех метрик")
	flag.Parse()

	return cfg
}

package tester

import (
	"burden/internal/config"
	"burden/internal/metrics"
)

func RunTest(cfg *config.Config) *metrics.Metrics {
	// TODO: main functional
	return metrics.CalculateMetrics()
}

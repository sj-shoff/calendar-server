package config

import (
	"flag"
	"os"
)

// Config представляет конфигурацию приложения
type Config struct {
	Port        string
	Environment string
}

// MustLoad загружает конфигурацию из переменных окружения и флагов
func MustLoad() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.Port, "port", "8888", "Port to run the server on")
	flag.StringVar(&cfg.Environment, "env", "development", "Application environment (development/production)")

	if port := os.Getenv("PORT"); port != "" {
		cfg.Port = port
	}
	if env := os.Getenv("ENVIRONMENT"); env != "" {
		cfg.Environment = env
	}

	flag.Parse()
	return cfg
}

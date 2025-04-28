package config

import (
	"flag"
)

var f ServerConfig

func parseFlags() *ServerConfig {
	if !flag.Parsed() {
		// Тест жестко требует порт 8080
		flag.StringVar(&f.ServerAddress, "a", ":8080", "server address")
		flag.StringVar(&f.BaseURL, "b", "http://localhost:8080", "base url")
		flag.Parse()
	}
	return &f
}

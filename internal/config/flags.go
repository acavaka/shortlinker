package config

import (
	"flag"
)

var f ServerConfig

func parseFlags() *ServerConfig {
	if !flag.Parsed() {
		flag.StringVar(&f.ServerAddress, "a", ":8080", "server address and port") // Слушаем все интерфейсы
		flag.StringVar(&f.BaseURL, "b", "http://[::1]:8080", "server address")    // IPv6-адрес
		flag.Parse()
	}
	return &f
}

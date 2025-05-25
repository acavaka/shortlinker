package config

import (
	"flag"
)

var f ServiceConfig

func parseFlags() *ServiceConfig {
	if !flag.Parsed() {
		flag.StringVar(&f.ServerAddress, "a", ":8080", "server address and port") // Слушаем все интерфейсы
		flag.StringVar(&f.BaseURL, "b", "http://[::1]:8080", "server address")    // IPv6-адрес
		flag.StringVar(&f.FileStoragePath, "f", "", "file path to save data")
		flag.Parse()
	}
	return &f
}

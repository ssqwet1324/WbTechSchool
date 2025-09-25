package flags

import (
	"flag"
	"log"
	"time"
)

// Flags - структура для хранения флагов
type Flags struct {
	Host    string
	Port    int
	Timeout time.Duration
}

// ParseFlags - парсим флаги
func ParseFlags() *Flags {
	host := flag.String("host", "localhost", "Host to connect")
	port := flag.Int("port", 8080, "Port to connect")
	timeout := flag.Duration("timeout", 10*time.Second, "Connection timeout")
	flag.Parse()

	if *host == "" || *port == 0 || *timeout == 0 {
		flag.PrintDefaults()
		log.Fatal("Host and Port are required")
	}

	return &Flags{*host, *port, *timeout}
}

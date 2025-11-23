package cli

import (
	"flag"
	"fmt"

	"L4.2/internal/entity"
)

// ParseServerFlags - парсим флаги для вооркеров
func ParseServerFlags() (entity.ServerFlags, string, string, error) {
	flagMode := flag.String("Mode", "leader", "режим запуска: leader или worker")
	flagAddr := flag.String("Addr", "localhost:8081", "адрес сервера")
	flagPeers := flag.String("Peers", "", "адреса воркеров через запятую")
	flagQuorum := flag.Int("Quorum", 1, "кворум ответивших серверов")
	flag.Parse()

	serverFlags := entity.ServerFlags{
		Mode:   *flagMode,
		Addr:   *flagAddr,
		Peers:  *flagPeers,
		Quorum: *flagQuorum,
	}

	args := flag.Args()
	if *flagMode == "worker" {
		return serverFlags, "", "", nil
	}

	if len(args) < 2 {
		return entity.ServerFlags{}, "", "", fmt.Errorf("для режима leader нужно указать <pattern> <filename>")
	}
	pattern := args[0]
	filename := args[1]
	return serverFlags, pattern, filename, nil
}

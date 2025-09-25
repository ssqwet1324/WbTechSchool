package main

import (
	"L2_17/internal/flags"
	"L2_17/internal/reader"
	"L2_17/internal/telnet"
	"log"
)

func main() {
	done := make(chan struct{})
	f := flags.ParseFlags()

	conn, err := telnet.ConnectTelnet(f)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := reader.Read(conn, done)
		if err != nil {
			log.Fatal("Ошибка чтения", err)
		}
	}()

	go func() {
		err := reader.Write(conn, done)
		if err != nil {
			log.Fatal("Ошибка ввода", err)
		}
	}()

	<-done
}

package reader

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
)

// Read - чтение данных из соединения и вывод в stdout
func Read(conn net.Conn, done chan struct{}) error {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println("Ответ сервера: ", scanner.Text())
	}

	close(done)
	return nil
}

// Write - чтение данных из stdin и отправка на сервер
func Write(conn net.Conn, done chan struct{}) error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		_, err := fmt.Fprintln(conn, scanner.Text())
		if err != nil {
			return fmt.Errorf("error writing line: %w", err)
		}
	}

	err := conn.Close()
	if err != nil {
		return errors.New("error closing connection" + err.Error())
	}
	close(done)
	return nil
}

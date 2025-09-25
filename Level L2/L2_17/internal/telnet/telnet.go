package telnet

import (
	"L2_17/internal/flags"
	"errors"
	"fmt"
	"net"
)

// ConnectTelnet - подключиться к telnet клиенту
func ConnectTelnet(flag *flags.Flags) (net.Conn, error) {
	addr := fmt.Sprintf("%s:%d", flag.Host, flag.Port)
	conn, err := net.DialTimeout("tcp", addr, flag.Timeout)
	if err != nil {
		return nil, errors.New("could not connect to telnet host: " + err.Error())
	}
	fmt.Println("Connect successful: ", conn.RemoteAddr())

	return conn, nil
}

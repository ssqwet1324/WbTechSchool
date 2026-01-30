package log

import (
	"fmt"
	"os"
	"sync"
)

// Log - структура логгера
type Log struct {
	LogChannel chan string
	startOnce  sync.Once
}

// New - конструктор
func New(logCh chan string) *Log {
	if logCh == nil {
		logCh = make(chan string, 100)
	}

	return &Log{
		LogChannel: logCh,
	}
}

// StartLogger - запускает горутину, которая читает канал и пишет в stdout.
func (l *Log) StartLogger() {
	if l == nil || l.LogChannel == nil {
		return
	}

	l.startOnce.Do(func() {
		go func() {
			for msg := range l.LogChannel {
				fmt.Fprintln(os.Stdout, msg)
			}
		}()
	})
}

// AsyncError - отправляет текст ошибки в асинхронный логгер
func (l *Log) AsyncError(msg string, err error) {
	if err == nil {
		return
	}
	l.enqueue(fmt.Sprintf("%s: %v", msg, err))
}

// AsyncMessage - отправляет произвольный текст в асинхронный логгер
func (l *Log) AsyncMessage(msg string) {
	l.enqueue(msg)
}

// AsyncMessagef - форматированый вывод
func (l *Log) AsyncMessagef(format string, args ...any) {
	l.enqueue(fmt.Sprintf(format, args...))
}

// enqueue - очередь для сообщений
func (l *Log) enqueue(msg string) {
	if l == nil || msg == "" {
		return
	}
	if l.LogChannel == nil {
		fmt.Fprintln(os.Stdout, msg)
		return
	}

	select {
	case l.LogChannel <- msg:
	default:
		// канал переполнен: чтобы не блокировать HTTP-хендлеры, выводим напрямую
		fmt.Fprintln(os.Stdout, msg)
	}
}

package queue

import (
	"errors"
)

// Queue - создаем очередь на отправку скачивания
func Queue(resources []string) (<-chan string, error) {
	if len(resources) == 0 {
		return nil, errors.New("queue: empty resources")
	}

	queue := make(chan string, len(resources))

	go func() {
		for _, r := range resources {
			queue <- r
		}
		close(queue)
	}()

	return queue, nil
}

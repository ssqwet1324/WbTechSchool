package service

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"

	"L4.2/internal/entity"
	"L4.2/internal/search"
)

// chunkResult — результат работы одного чанка
type chunkResult struct {
	index int
	lines []string
	count int
	err   error
}

// Execute запускает локальный или распределённый поиск
func Execute(lines []string, peers []string, options entity.Options, pattern string, quorum int) error {
	// Если воркеров нет — работаем локально
	if len(peers) == 0 {
		res, err := search.SearchLines(lines, options, pattern, 0)
		if err != nil {
			return err
		}
		outputResult(res, options)
		return nil
	}

	// Делим входные строки
	chunks := splitIntoChunks(lines, len(peers)+1)

	// Для каждого чанка считаем смещение (чтобы знать его позицию в исходном файле)
	offsets := calculateOffsets(chunks)

	// Запускаем обработку чанков
	results := processChunks(chunks, offsets, peers, options, pattern)

	// Проверяем кворум — достаточно ли успешных ответов
	if len(results) < quorum {
		return fmt.Errorf("не достигнут кворум: получено %d, требуется %d", len(results), quorum)
	}

	// Выводим объединённый результат
	outputResults(results, options)

	return nil
}

// processChunks параллельно выполняет поиск для каждого чанка
func processChunks(chunks [][]string, offsets []int, peers []string, options entity.Options, pattern string) []chunkResult {
	var wg sync.WaitGroup
	results := make(chan chunkResult, len(chunks))

	// Запускаем обработку каждого чанка
	for i, chunk := range chunks {
		if len(chunk) == 0 {
			continue
		}

		wg.Add(1)
		go func(i int, chunk []string) {
			defer wg.Done()

			var (
				res search.SearchResult
				err error
			)

			if i == 0 {
				res, err = search.SearchLines(chunk, options, pattern, offsets[i])
			} else {
				peer := peers[(i-1)%len(peers)]
				res, err = dispatchToPeer(peer, pattern, chunk, options, offsets[i])
			}

			// Возвращаем результат
			if err != nil {
				results <- chunkResult{index: i, err: err}
			} else {
				results <- chunkResult{index: i, lines: res.Lines, count: res.Count}
			}
		}(i, chunk)
	}

	// Закрываем канал после завершения всех горутин
	go func() {
		wg.Wait()
		close(results)
	}()

	// Читаем все результаты
	var all []chunkResult
	for r := range results {
		if r.err != nil {
			fmt.Printf("Ошибка chunk %d: %v\n", r.index, r.err)
			continue
		}
		all = append(all, r)
	}

	return all
}

// outputResults выводит объединённый результат распределённого поиска
func outputResults(results []chunkResult, options entity.Options) {
	// Сортировка по индексу чанка (чтобы сохранить реальный порядок строк)
	sort.Slice(results, func(i, j int) bool {
		return results[i].index < results[j].index
	})

	// Если включён режим -c (count) — суммируем количество совпадений
	if options.Count {
		total := 0
		for _, r := range results {
			total += r.count
		}
		fmt.Println(total)
		return
	}

	// Вывод всех найденных строк
	for _, r := range results {
		for _, ln := range r.lines {
			fmt.Println(ln)
		}
	}
}

// outputResult выводит результат локального поиска
func outputResult(res search.SearchResult, options entity.Options) {
	if options.Count {
		fmt.Println(res.Count)
		return
	}
	for _, ln := range res.Lines {
		fmt.Println(ln)
	}
}

// splitIntoChunks делит строки на N чанков
func splitIntoChunks(lines []string, num int) [][]string {
	if num < 1 {
		num = 1
	}

	// Размер одного чанка
	size := (len(lines) + num - 1) / num

	chunks := make([][]string, num)
	for i := 0; i < num; i++ {
		start := i * size
		if start >= len(lines) {
			break
		}
		end := start + size
		if end > len(lines) {
			end = len(lines)
		}
		chunks[i] = lines[start:end]
	}

	return chunks
}

// calculateOffsets считает, сколько строк было до каждого чанка
func calculateOffsets(chunks [][]string) []int {
	offsets := make([]int, len(chunks))
	sum := 0
	for i, c := range chunks {
		offsets[i] = sum
		sum += len(c)
	}
	return offsets
}

// dispatchToPeer отправляет запрос к удалённому воркеру
func dispatchToPeer(peer, pattern string, lines []string, options entity.Options, offset int) (search.SearchResult, error) {
	// Формируем JSON запрос
	body, _ := json.Marshal(entity.SearchRequest{
		Pattern: pattern,
		Lines:   lines,
		Options: options,
		Offset:  offset,
	})

	resp, err := http.Post("http://"+peer+"/search", "application/json", bytes.NewReader(body))
	if err != nil {
		return search.SearchResult{}, fmt.Errorf("запрос к %s: %w", peer, err)
	}
	defer resp.Body.Close()

	// Проверяем код ответа
	if resp.StatusCode != http.StatusOK {
		return search.SearchResult{}, fmt.Errorf("код %d от %s", resp.StatusCode, peer)
	}

	var payload entity.SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return search.SearchResult{}, fmt.Errorf("decode %s: %w", peer, err)
	}

	return search.SearchResult{
		Lines: payload.Lines,
		Count: payload.Count,
	}, nil
}

// ReadAll читает все строки
func ReadAll(scanner *bufio.Scanner) ([]string, error) {
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// ParsePeers разбирает список адресов
func ParsePeers(s string) []string {
	if s == "" {
		return nil
	}

	parts := strings.Split(s, ",")
	res := make([]string, 0, len(parts))

	for _, p := range parts {
		if addr := strings.TrimSpace(p); addr != "" {
			res = append(res, addr)
		}
	}

	return res
}

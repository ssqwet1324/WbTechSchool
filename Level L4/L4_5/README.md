# Api Optimization

## Используемые технологии

- **Go 1.25.3** - язык программирования
- **Gin** - веб-фреймворк для HTTP API
- **Pprof** - профилировщик

## RPS ручки сервера без оптимизации

**Команда для нагрузки в пике**
```bash
hey -n 10000 -c 50 -m POST -H "Content-Type: application/json" -d '{"a":1,"b":2}' http://localhost:8080/sum
```
**Результат выполнения**
```
Summary:
  Total:        0.1531 secs
  Slowest:      0.0218 secs
  Fastest:      0.0000 secs
  Average:      0.0007 secs
  Requests/sec: 65297.8070

  Total data:   90000 bytes
  Size/request: 9 bytes
```

**Средняя производительность**
```bash
hey -z 40s -c 50 -m POST ` -H "Content-Type: application/json" ` -d "{""a"":1,""b"":2}" http://localhost:8080/sum
```
**Результат выполнения**
```
Summary:
  Total:        40.0017 secs
  Slowest:      0.0190 secs
  Fastest:      0.0000 secs
  Average:      0.0020 secs
  Requests/sec: 92034.8150

  Total data:   33133977 bytes
  Size/request: 33 bytes
```

**Данные профилировщика:**
**Cpu Profile**
```
Showing top 10 nodes out of 173
     flat  flat%   sum%        cum   cum%
     60.57s 24.28% 24.28%     62.75s 25.16%  runtime.cgocall
     56.48s 22.64% 46.93%     56.52s 22.66%  runtime.stdcall0
     43.63s 17.49% 64.42%     43.63s 17.49%  runtime.stdcall2
     4.34s  1.74% 66.16%      4.34s  1.74%  runtime.stdcall1
     3.66s  1.47% 67.62%      3.66s  1.47%  runtime.stdcall6
     1.97s  0.79% 68.41%      1.97s  0.79%  runtime.procyield
     1.44s  0.58% 68.99%      1.44s  0.58%  runtime.nextFreeFast
     1.34s  0.54% 69.53%      1.79s  0.72%  runtime.findObject
     0.82s  0.33% 69.86%      2.54s  1.02%  runtime.scanobject
     0.68s  0.27% 70.13%      1.53s  0.61%  internal/runtime/maps.(*Map).getWithoutKeySmallFastStr
```
**Heap Profile**
```
Showing nodes accounting for 7715.42kB, 100% of 7715.42kB total
Showing top 10 nodes out of 71
    flat  flat%   sum%        cum   cum%
    1539kB 19.95% 19.95%     1539kB 19.95%  runtime.allocm
    1050.86kB 13.62% 33.57%  1050.86kB 13.62%  github.com/go-playground/validator/v10.map.init.7
    1028kB 13.32% 46.89%     1028kB 13.32%  bufio.NewWriterSize
    1024.62kB 13.28% 60.17%  1024.62kB 13.28%  net.newFD (inline)
    1024.34kB 13.28% 73.45%  1024.34kB 13.28%  github.com/gin-gonic/gin/render.writeContentType
    512.25kB  6.64% 80.09%   512.25kB  6.64%  encoding/json.(*Decoder).refill
    512.16kB  6.64% 86.73%   512.16kB  6.64%  net/http.readRequest
    512.14kB  6.64% 93.36%   512.14kB  6.64%  reflect.addReflectOff
    512.05kB  6.64%   100%   512.05kB  6.64%  internal/poll.runtime_Semacquire
    0     0%   100%  1024.62kB 13.28%  api_optimization/internal/app.Run
```

**Данные из бенчмарка**
```
goos: linux
goarch: amd64
pkg: api_optimization/internal/handler
cpu: 12th Gen Intel(R) Core(TM) i5-12500H
BenchmarkSumHandler-16    	  285045	      3588 ns/op	    8594 B/op	      40 allocs/op
BenchmarkSumHandler-16    	  339470	      3582 ns/op	    8593 B/op	      40 allocs/op
BenchmarkSumHandler-16    	  334414	      3521 ns/op	    8593 B/op	      40 allocs/op
BenchmarkSumHandler-16    	  333220	      3458 ns/op	    8594 B/op	      40 allocs/op
BenchmarkSumHandler-16    	  345970	      3452 ns/op	    8594 B/op	      40 allocs/op
BenchmarkSumHandler-16    	  334849	      3465 ns/op	    8594 B/op	      40 allocs/op
BenchmarkSumHandler-16    	  341671	      3515 ns/op	    8594 B/op	      40 allocs/op
BenchmarkSumHandler-16    	  347898	      3446 ns/op	    8594 B/op	      40 allocs/op
BenchmarkSumHandler-16    	  349864	      3459 ns/op	    8594 B/op	      40 allocs/op
BenchmarkSumHandler-16    	  346146	      3537 ns/op	    8594 B/op	      40 allocs/op
PASS
ok  	api_optimization/internal/handler	12.158s
```

## RPS ручки сервера после оптимизаций

**Команда для нагрузки в пике**
```bash
hey -n 10000 -c 50 -m POST -H "Content-Type: application/json" -d '{"a":1,"b":2}' http://localhost:8080/sum
```
**Результат выполнения**
```
Summary:
  Total:        0.1388 secs
  Slowest:      0.0177 secs
  Fastest:      0.0000 secs
  Average:      0.0007 secs
  Requests/sec: 72046.0057

  Total data:   90000 bytes
  Size/request: 9 bytes
```

**Средняя производительность**
```bash
hey -z 40s -c 50 -m POST ` -H "Content-Type: application/json" ` -d "{""a"":1,""b"":2}" http://localhost:8080/sum
```
**Результат выполнения**
```
Summary:
  Total:        40.0016 secs
  Slowest:      0.0244 secs
  Fastest:      0.0000 secs
  Average:      0.0020 secs
  Requests/sec: 105959.9046

  Total data:   38147085 bytes
  Size/request: 38 bytes
```

**Данные профилировщика:**
**Cpu Profile**
```
Showing top 10 nodes out of 167
      flat  flat%   sum%        cum   cum%
     69.28s 26.72% 26.72%     71.53s 27.59%  runtime.cgocall
     63.69s 24.57% 51.29%     63.70s 24.57%  runtime.stdcall0
     37.73s 14.55% 65.84%     37.73s 14.55%  runtime.stdcall2
     4.18s  1.61% 67.46%      4.18s  1.61%  runtime.stdcall1
     3.92s  1.51% 68.97%      3.93s  1.52%  runtime.stdcall6
     1.92s  0.74% 69.71%      1.92s  0.74%  runtime.procyield
     1.75s  0.68% 70.38%      1.75s  0.68%  runtime.nextFreeFast
     1.41s  0.54% 70.93%      1.85s  0.71%  runtime.findObject
     1.02s  0.39% 71.32%      2.53s  0.98%  runtime.scanobject
     0.96s  0.37% 71.69%      2.92s  1.13%  time.Time.appendFormat
```
**Heap Profile**
```
Showing nodes accounting for 2574.05kB, 100% of 2574.05kB total
Showing top 10 nodes out of 18
      flat  flat%   sum%        cum   cum%
      1539kB 59.79% 59.79%     1539kB 59.79%  runtime.allocm
      521.05kB 20.24% 80.03%   521.05kB 20.24%  encoding/xml.map.init.0
      514kB 19.97%   100%      514kB 19.97%  bufio.NewReaderSize (inline)
      0     0%   100%      514kB 19.97%  bufio.NewReader (inline)
      0     0%   100%   521.05kB 20.24%  encoding/xml.init
      0     0%   100%      514kB 19.97%  net/http.(*conn).serve
      0     0%   100%      514kB 19.97%  net/http.newBufioReader
      0     0%   100%   521.05kB 20.24%  runtime.doInit (inline)
      0     0%   100%   521.05kB 20.24%  runtime.doInit1
      0     0%   100%   521.05kB 20.24%  runtime.main
```

**Данные из бенчмарка**
```
goos: linux
goarch: amd64
pkg: api_optimization/internal/handler
cpu: 12th Gen Intel(R) Core(TM) i5-12500H
BenchmarkSumHandler-16    	  325923	      3467 ns/op	    8701 B/op	      37 allocs/op
BenchmarkSumHandler-16    	  338008	      3656 ns/op	    8718 B/op	      37 allocs/op
BenchmarkSumHandler-16    	  313965	      3358 ns/op	    8688 B/op	      37 allocs/op
BenchmarkSumHandler-16    	  346648	      3403 ns/op	    8689 B/op	      37 allocs/op
BenchmarkSumHandler-16    	  364416	      3352 ns/op	    8718 B/op	      37 allocs/op
BenchmarkSumHandler-16    	  363804	      3272 ns/op	    8682 B/op	      37 allocs/op
BenchmarkSumHandler-16    	  379890	      3280 ns/op	    8698 B/op	      37 allocs/op
BenchmarkSumHandler-16    	  380512	      3409 ns/op	    8705 B/op	      37 allocs/op
BenchmarkSumHandler-16    	  343026	      3348 ns/op	    8684 B/op	      37 allocs/op
BenchmarkSumHandler-16    	  359361	      3301 ns/op	    8717 B/op	      37 allocs/op
PASS
ok  	api_optimization/internal/handler	12.282s
```

## Сравнение данных из бенчмарков
```
goos: linux
goarch: amd64
pkg: api_optimization/internal/handler
cpu: 12th Gen Intel(R) Core(TM) i5-12500H
              │   old.txt   │              new.txt               │
              │   sec/op    │   sec/op     vs base               │
SumHandler-16   3.490µ ± 3%   3.355µ ± 3%  -3.87% (p=0.007 n=10)

              │   old.txt    │               new.txt               │
              │     B/op     │     B/op      vs base               │
SumHandler-16   8.393Ki ± 0%   8.496Ki ± 0%  +1.23% (p=0.000 n=10)

              │  old.txt   │              new.txt              │
              │ allocs/op  │ allocs/op   vs base               │
SumHandler-16   40.00 ± 0%   37.00 ± 0%  -7.50% (p=0.000 n=10)
```

**Итоговый прирост**
```
| Метрика | Исходный (old.txt) | Оптимизированный (new.txt) | Прирост                |
| ------- | ------------------ | -------------------------- | ---------------------- |
| CPU     | 3.490µs ± 3%       | 3.355µs ± 3%               | -3.87% (быстрее)       |
| Память  | 8.393KiB           | 8.496KiB                   | +1.23% (незначительно) |
| Allocs  | 40                 | 37                         | -7.50% (прирост ~8%)   |
```

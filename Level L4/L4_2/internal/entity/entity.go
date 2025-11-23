package entity

// Options - структура для хранения данных флагов
type Options struct {
	After   int  // -A
	Before  int  // -B
	Count   bool // -c
	Ignore  bool // -i
	Invert  bool // -v
	Fixed   bool // -F
	LineNum bool // -n
}

// ServerFlags - структура для флагов сервиса
type ServerFlags struct {
	Mode   string
	Addr   string
	Peers  string
	Quorum int
}

type SearchRequest struct {
	Pattern string   `json:"pattern"`
	Lines   []string `json:"lines"`
	Options Options  `json:"options"`
	Offset  int      `json:"offset"`
}

type SearchResponse struct {
	Lines []string `json:"lines"`
	Count int      `json:"count"` // количество совпадений (для -c)
}

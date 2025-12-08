package entity

// Input - структура для ввода чисел
type Input struct {
	A int `json:"a"`
	B int `json:"b"`
}

// Output - вывод суммы
type Output struct {
	Sum int `json:"sum"`
}

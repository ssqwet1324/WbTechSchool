package main

import "fmt"

type Human struct {
	Name  string
	Phone string
}

func (h *Human) SayHi() {
	fmt.Printf("Hi, I am %s, my phone %s\n", h.Name, h.Phone)
}

func (h *Human) SayGoodbye() {
	fmt.Printf("Goodbye %s, my phone %s\n", h.Name, h.Phone)
}

type Action struct {
	Age string
	Human
}

func main() {
	//тут инициализируем структуру
	a := Action{Age: "13", Human: Human{"Gik", "123"}}
	//тут вызываем метод из Human благодаря встраиванию
	a.SayHi()
	a.SayGoodbye()
	//тут выводим поля благодаря встраиванию
	fmt.Println(a.Age, a.Name, a.Phone)
	//Выводит:
	//Hi, I am Gik, my phone 123
	//Goodbye Gik, my phone 123
	//13 Gik 123

}

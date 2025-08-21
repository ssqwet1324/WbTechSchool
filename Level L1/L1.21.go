package main

import "fmt"

type NewMethod interface {
	New()
}

type OldMethod struct{}

func (o *OldMethod) Old() {
	fmt.Println("Я тут")
}

type Adapter struct {
	OldMethod *OldMethod
}

func (a *Adapter) New() {
	a.OldMethod.Old()
}

func main() {
	var n NewMethod = &Adapter{OldMethod: &OldMethod{}}
	n.New()
}

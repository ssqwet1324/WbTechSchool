package main

import (
	"fmt"
	"math"
)

type Point struct {
	x float64
	y float64
}

func (p *Point) Distance(other *Point) float64 {
	d := math.Sqrt(math.Pow(other.x-p.x, 2) + math.Pow(other.y-p.y, 2))
	return d
}

func NewPoint(x, y float64) *Point {
	return &Point{x, y}
}

func main() {
	var x1, y1, x2, y2 float64
	if _, err := fmt.Scan(&x1, &y1, &x2, &y2); err != nil {
		panic("Введено не число с плавающей запятой")
	}
	p1 := NewPoint(x1, y1)
	p2 := NewPoint(x2, y2)
	fmt.Println(p1.Distance(p2))
}

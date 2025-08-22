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
	p1 := NewPoint(1.0, 2.0)
	p2 := NewPoint(2.0, 3.0)
	fmt.Println(p1.Distance(p2))
}

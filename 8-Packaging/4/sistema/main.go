package main

import "github.com/rafael0502/curso-go/8-Packaging/4/math"

func main() {
	a := math.NewMath(1, 2)
	println(a.Add())
}

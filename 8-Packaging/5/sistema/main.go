package main

//go work init ./math ./sistema
//go mod tidy -e

import (
	"github.com/google/uuid"
	"github.com/rafael0502/curso-go/8-Packaging/4/math"
)

func main() {
	a := math.NewMath(1, 2)
	println(a.Add())
	println(uuid.New().String())
}

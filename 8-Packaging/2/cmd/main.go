package main

import (
	"fmt"

	"github.com/rafael0502/curso-go/8-Packaging/2/math"
)

func main() {
	m := math.NewMath(5, 3)

	fmt.Println(m.Add())
	fmt.Println(m)
}

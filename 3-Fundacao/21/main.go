package main

import (
	"fmt"

	"curso-go/matematica"
)

func main() {
	s := matematica.Soma(10, 5)
	fmt.Println("Resultado: ", s)
	fmt.Println(matematica.A)
}

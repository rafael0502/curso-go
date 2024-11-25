package main

import (
	"errors"
	"fmt"
)

func main() {
	valor, err := sum(10, 6)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(valor)

	//fmt.Println(sum(50, 6))
}

func sum(a, b int) (int, error) {
	if a+b >= 50 {
		return 0, errors.New("A soma é maior/igual a 50")
	}
	return a + b, nil
}

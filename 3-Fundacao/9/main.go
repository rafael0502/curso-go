package main

import "fmt"

func main() {
	fmt.Print(sum(1, 3, 4, 5, 55, 70))
}

func sum(numeros ...int) int {
	total := 0

	for _, numero := range numeros {
		total += numero
	}

	return total
}

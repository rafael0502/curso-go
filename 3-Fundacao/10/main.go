package main

import "fmt"

func main() {
	total := func() int {
		return sum(1, 3, 4, 5, 55, 70) * 2
	}()

	fmt.Print(total)
}

func sum(numeros ...int) int {
	total := 0

	for _, numero := range numeros {
		total += numero
	}

	return total
}

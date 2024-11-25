package main

import "fmt"

func main() {
	salarios := map[string]int{"Rafael": 8000, "João": 2000, "Maria": 5000}
	//	fmt.Println(salarios["Rafael"])

	for nome, salario := range salarios {
		fmt.Printf("O salário de %s é %d\n", nome, salario)
	}

	for _, salario := range salarios {
		fmt.Printf("O salário é %d\n", salario)
	}
}

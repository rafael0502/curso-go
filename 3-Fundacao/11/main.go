package main

import "fmt"

type Cliente struct {
	Nome  string
	Idade int
	Ativo bool
}

func main() {
	cli := Cliente{
		Nome:  "Rafael",
		Idade: 42,
		Ativo: true,
	}

	cli.Ativo = false

	fmt.Println("Nome: " + cli.Nome)
	fmt.Println(cli.Ativo)
}

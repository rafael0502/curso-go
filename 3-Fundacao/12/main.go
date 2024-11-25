package main

import "fmt"

type Endereco struct {
	Logradouro string
	Numero     int
	Cidade     string
	Estado     string
}

type Cliente struct {
	Nome  string
	Idade int
	Ativo bool
	Endereco
}

func main() {
	cli := Cliente{
		Nome:  "Rafael",
		Idade: 42,
		Ativo: true,
	}

	cli.Ativo = false
	//	cli.Cidade = "Palhoça"  // ou
	cli.Endereco.Cidade = "Palhoça"

	fmt.Println("Nome: " + cli.Nome)
	fmt.Println(cli.Ativo)
}

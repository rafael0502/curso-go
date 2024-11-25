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

func (c Cliente) Desativar() {
	c.Ativo = false
	fmt.Println("O cliente foi desativado")
}

func main() {
	cli := Cliente{
		Nome:  "Rafael",
		Idade: 42,
		Ativo: true,
	}

	cli.Desativar()
}

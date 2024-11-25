package main

import "fmt"

type Endereco struct {
	Logradouro string
	Numero     int
	Cidade     string
	Estado     string
}

type Pessoa interface {
	Desativar()
}

type Empresa struct {
	Nome string
}

func (e Empresa) Desativar() {

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

func Desativacao(pessoa Pessoa) {
	pessoa.Desativar()
}

func main() {
	cli := Cliente{
		Nome:  "Rafael",
		Idade: 42,
		Ativo: true,
	}

	minhaEmp := Empresa{}

	Desativacao(cli)
	Desativacao(minhaEmp)
}

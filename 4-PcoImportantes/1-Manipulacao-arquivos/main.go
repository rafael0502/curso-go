package main

import (
	"os"
)

func main() {
	/*	f, err := os.Create("arquivo.txt")

		if err != nil {
			panic(err)
		}

		//tamanho, err := f.WriteString("Dae cambada!")
		tamanho, err := f.Write([]byte("Escrevendo dados no arquivo..."))

		if err != nil {
			panic(err)
		}

		fmt.Printf("Arquivo criado com sucesso! Tamanho: %d bytes", tamanho)
		f.Close()
	*/
	//Leitura...
	/*
		arquivo, err := os.ReadFile("arquivo.txt")

		if err != nil {
			panic(err)
		}

		fmt.Println(string(arquivo))
	*/
	/*
		arq, err := os.Open("arquivo.txt")

		if err != nil {
			panic(err)
		}

		reader := bufio.NewReader(arq)
		buffer := make([]byte, 3)

		for {
			n, err := reader.Read(buffer)

			if err != nil {
				break
			}

			fmt.Println(string(buffer[:n]))
		}
	*/
	err := os.Remove("arquivo.txt")

	if err != nil {
		panic(err)
	}
}

package main

import (
	"fmt"
)

// canal que somente recebe chan<-
func recebe(nome string, hello chan<- string) {
	hello <- nome
}

// canal que somente envia <-chan
func ler(data <-chan string) {
	fmt.Println(<-data)
}

func main() {
	hello := make(chan string)
	go recebe("Hello", hello)
	ler(hello)
}

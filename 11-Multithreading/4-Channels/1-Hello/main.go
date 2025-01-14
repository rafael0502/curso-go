package main

import "fmt"

// Thread 1
func main() {
	// Create a channel
	canal := make(chan string) //vazio

	// Thread 2
	go func() {
		canal <- "Olá mundo!" // Está cheio
	}()

	// Thread 1
	msg := <-canal // canal esvaziado
	fmt.Println(msg)
}

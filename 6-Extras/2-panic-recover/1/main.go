package main

import "fmt"

func panic1() {
	panic("Panico 1")
}

func panic2() {
	panic("Panico 2")
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			if r == "Panico 1" {
				fmt.Println("Recuperou o panic 1")
			}

			if r == "Panico 2" {
				fmt.Println("Recuperou o panic 2")
			}
		}
	}()

	panic1()
}

package main

func main() {
	a := 10
	var ponteiro *int = &a // * indica que é um ponteiro
	*ponteiro = 20
	b := &a
	*b = 30
	//println(b)
	println(*b)
}

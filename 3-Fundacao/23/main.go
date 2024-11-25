package main

func main() {
	a := 4
	b := 2
	c := 3

	// condição 1 e  cond 2: a>b && c>a
	// condição 1 ou cond 2: a>b || c>a

	if a > b && c > a {
		println(a)
	} else {
		println(b)
	}

	switch a {
	case 1:
		println("a")
	case 2:
		println("b")
	case 3:
		println("c")
	default:
		println("d")
	}
}

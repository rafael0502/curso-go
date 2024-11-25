package main

//import "fmt"
import (
	"fmt"
)

const a = "Dae mundo!!"

type ID int

var (
	b bool    = true
	c int     = 1
	d string  = "Rafa"
	e float64 = 1.9
	f ID      = 1
)

func main() {
	var arrTeste [3]int

	arrTeste[0] = 31
	arrTeste[1] = 22
	arrTeste[2] = 31

	//fmt.Printf("O tipo de E é %T", f)
	for i, v := range arrTeste {
		fmt.Printf("O valor do indice é %d e o valor é %d\n", i, v)
	}
}

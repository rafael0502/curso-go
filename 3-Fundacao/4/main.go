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
	fmt.Printf("O tipo de E é %T", f)
}

package main

import (
	"fmt"

	"github.com/rafael0502/utils-secret/pkg/events"
)

func main() {
	ed := events.NewEventDispatcher()
	fmt.Println(ed)
}

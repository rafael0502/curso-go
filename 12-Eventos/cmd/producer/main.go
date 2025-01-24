package main

import "github.com/rafael0502/curso-go/utils/pkg/rabbitmq"

func main() {
	ch, err := rabbitmq.OpenChannel()

	if err != nil {
		panic(err)
	}

	defer ch.Close()

	rabbitmq.Publish(ch, "Dae raça!!", "amq.direct")
}

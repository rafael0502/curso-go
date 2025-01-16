package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

type Message struct {
	id  int64
	Msg string
}

func main() {
	c1 := make(chan Message)
	c2 := make(chan Message)
	var i int64 = 0

	// RabbitMQ
	go func() {
		for {
			time.Sleep(2 * time.Second)
			atomic.AddInt64(&i, 1) // atomic operation to increment i para que não ocorra concorrência
			msg := Message{id: i, Msg: "Hello from RabbitMQ"}
			c1 <- msg
		}
	}()

	// Kafka
	go func() {
		for {
			time.Sleep(1 * time.Second)
			atomic.AddInt64(&i, 1)
			msg := Message{id: i, Msg: "Hello from Kafka"}
			c2 <- msg
		}
	}()

	//for i := 0; i < 3; i++ {
	for {
		select {
		case msg := <-c1: //RabbitMQ
			fmt.Printf("Received from RabbitMQ: ID: %d - %s\n", msg.id, msg.Msg)
		case msg := <-c2: //Kafka
			fmt.Printf("Received from Kafka: ID: %d - %s\n", msg.id, msg.Msg)
		case <-time.After(time.Second * 3):
			println("timeout")
			//default:
			//	println("default")
		}
	}
}

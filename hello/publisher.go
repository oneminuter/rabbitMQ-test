package main

import (
	"log"
	"rabbitmq_test/config"

	"github.com/streadway/amqp"
)

/*
go run publisher.go
*/
func main() {
	conn, err := amqp.Dial(config.DailUrl)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to connect to RabbitMQ", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "Failed to open ja channel", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare a queue", err)
	}

	body := "Hello World111"
	err = ch.Publish("", q.Name, false, false, amqp.Publishing{
		ContentType: "text/palin",
		Body:        []byte(body),
	})
	if err != nil {
		log.Fatalf("%s: %s", "Failed to publish a message", err)
	}
}

package main

import (
	"log"
	"rabbitmq_test/config"

	"github.com/streadway/amqp"
)

/*
	go run consumer.go
*/
func main() {
	conn, err := amqp.Dial(config.DailUrl)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to connect to RabbitMQ", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "Failed to open a channel", err)
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
	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to register a consumer", err)
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Printf(" [*] Waiting for message. To exit press CTRL+C")
	<-forever
}

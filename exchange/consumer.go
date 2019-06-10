package main

import (
	"log"
	"rabbitmq_test/config"

	"github.com/streadway/amqp"
)

/*
	Publish/Subscribe

	Listing exchanges:
		rabbitmqctl list_exchanges
	Listing bindings:
	rabbitmqctl list_bindings

	go run consumer.go > logs_from_rabbit.log
	go run consumer.go
*/
/*
	broadcast
*/

func main() {
	conn, err := amqp.Dial(config.DailUrl)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("%s: $s", "Failed to open a channel", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare an exchange", err)
	}

	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare a queue", err)
	}

	err = ch.QueueBind(q.Name, "", "logs", false, nil)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to bind a queue", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("%s: %s", "failed to register a consumer", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press ctrl+c")
	<-forever
}

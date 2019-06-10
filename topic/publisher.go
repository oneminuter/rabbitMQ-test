package main

import (
	"log"
	"os"
	"rabbitmq_test/config"
	"rabbitmq_test/util"

	"github.com/streadway/amqp"
)

/*
go run publisher.go "kern.critical" "A critical kernel error"
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

	err = ch.ExchangeDeclare(
		"logs_topic",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("%s： %s", "Failed to declare an exchange", err)
	}

	body := util.BodyFrom(os.Args)
	err = ch.Publish("logs_topic", util.SeverityFrom(os.Args), false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to publish a message", err)
	}
	log.Printf(" [x] Send: %s", body)
}

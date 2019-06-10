package main

import (
	"log"
	"os"
	"rabbitmq_test/config"
	"rabbitmq_test/util"

	"github.com/streadway/amqp"
)

/*
	go run publisher.go "this is messages"
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

	err = ch.ExchangeDeclare("logs", "fanout", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("%s: $s", "Failed to declare exchange", err)
	}

	body := util.BodyFrom(os.Args)
	err = ch.Publish("logs", "", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	})
	if err != nil {
		log.Fatalf("%s: $s", "Failed to publish ", err)
	}

}

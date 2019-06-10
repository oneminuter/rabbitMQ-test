package main

import (
	"log"
	"os"
	"rabbitmq_test/config"
	"rabbitmq_test/util"

	"github.com/streadway/amqp"
)

/*
go run publisher.go First message.
go run publisher.go Second message..
go run publisher.go Third message...
go run publisher.go Fourth message....
go run publisher.go Fifth message.....
...
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
		"task_queue", //it won't work in our present setup. That's because we've already defined a queue called hello which is not durable. RabbitMQ doesn't allow you to redefine an existing queue with different parameters and will return an error to any program that tries to do that
		true,         //When RabbitMQ quits or crashes it will forget the queues and messages unless you tell it not to. Two things are required to make sure that messages aren't lost: we need to mark both the queue and messages as durable.
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare a queue", err)
	}

	body := util.BodyFrom(os.Args)
	err = ch.Publish("", q.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent, //Marking messages as persistent
		ContentType:  "text/plain",
		Body:         []byte(body),
	})
	if err != nil {
		log.Fatalf("%s: %s", "Failed to publish a message", err)
	}
}

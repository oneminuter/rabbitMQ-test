package main

import (
	"log"
	"os"
	"rabbitmq_test/config"

	"github.com/streadway/amqp"
)

/*
When a queue is bound with "#" (hash) binding key - it will receive all the messages, regardless of the routing key - like in fanout exchange.
When special characters "*" (star) and "#" (hash) aren't used in bindings, the topic exchange will behave just like a direct one.

go run consumer.go "#"
go run consumer.go "kern.*"
go run consumer.go "*.critical"
go run consumer.go "kern.*" "*.critical"
*/
/*
	每个订阅者都能收到对应 topic 下发的消息
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
		"logs_tpoic",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare an exhange", err)
	}

	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to dclare a queue")
	}

	if len(os.Args) < 2 {
		log.Printf("Usage: %s []binding_key]...", os.Args[0])
		os.Exit(0)
	}
	for _, s := range os.Args[1:] {
		log.Printf("Binding queue %s to exchange %s with routing eu %s", q.Name, "logs_topic", s)
		err = ch.QueueBind(q.Name, s, "logs_topic", false, nil)
		if err != nil {
			log.Fatalf("%s: %s", "Failed to bind a queue", err)
		}
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to register a consumer", err)
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

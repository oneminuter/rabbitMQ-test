package main

import (
	"log"
	"os"
	"rabbitmq_test/config"

	"github.com/streadway/amqp"
)

/*
	go run consumer.go warning error > logs_from_rabbit.log
	go run consumer.go info warning error
*/

/*
	按指定消息特征接收，应用场景，比如接收不同的日志级别进行处理
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
		"logs_direct",
		"direct",
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

	if len(os.Args) < 2 {
		log.Printf("Usage: %s [info] [warning] [error]", os.Args[0])
		os.Exit(0)
	}

	for _, s := range os.Args[1:] {
		log.Printf("Binding queue %s to exchange %s with routing key %s", q.Name, "logs_direct", s)
		err = ch.QueueBind(q.Name, s, "logs_direct", false, nil)
		if err != nil {
			log.Fatalf("%s: %s", "Failed to bind a queue", err)
		}
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

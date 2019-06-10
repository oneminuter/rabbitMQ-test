package main

import (
	"bytes"
	"log"
	"rabbitmq_test/config"
	"time"

	"github.com/streadway/amqp"
)

/*
go run consumer.go
go run consumer.go
...
*/
/*
	消息会平均到每个消费者
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
		"task_queue",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare a queue", err)
	}

	err = ch.Qos(1, 0, false)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to set Qos", err)
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to register a consumer", err)
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			dot_count := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dot_count)
			time.Sleep(t * time.Second)
			log.Printf("done")
			d.Ack(false)
		}
	}()
	log.Printf(" [*] Waiting for message. To exit press CTRL+C")
	<-forever
}

package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
)

var rabbit_host = "localhost"
var rabbit_port = "5672"
var rabbit_user = "guest"
var rabbit_password = "guest"

func main() {
	if len(os.Args) >= 2 {
		rabbit_host = os.Args[1]
	}
	consume()
}

func consume() {
	var rabbit_url = "amqp://" + rabbit_user + ":" + rabbit_password + "@" + rabbit_host + ":" + rabbit_port + "/"
	conn, err := amqp.Dial(rabbit_url)

	if err != nil {
		log.Fatalf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	ch, err := conn.Channel()

	if err != nil {
		log.Fatalf("%s: %s", "Failed to open a channel", err)
	}

	q, err := ch.QueueDeclare(
		"publisher", // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)

	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare a queue", err)
	}

	fmt.Println("Channel and Queue established")

	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		log.Fatalf("%s: %s", "Failed to register consumer", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			d.Ack(false)
		}
	}()

	fmt.Println("Running...")
	<-forever
}

package main

import (
	"fmt"
	"os"
	"time"

	// "github.com/streadway/amqp"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	// "os"
)

// var rabbit_host = os.Getenv("RABBIT_HOST")
// var rabbit_port = os.Getenv("RABBIT_PORT")
// var rabbit_user = os.Getenv("RABBIT_USERNAME")
// var rabbit_password = os.Getenv("RABBIT_PASSWORD")

var rabbit_host = "localhost"
var rabbit_port = "5672"
var rabbit_user = "guest"
var rabbit_password = "guest"

func main() {
	if len(os.Args) >= 2 {
		rabbit_host = os.Args[1]
	}
	var rabbit_url = "amqp://" + rabbit_user + ":" + rabbit_password + "@" + rabbit_host + ":" + rabbit_port + "/"
	var message = "hello world to " + rabbit_host

	conn, err := amqp.Dial(rabbit_url)

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

	// ticker := time.NewTicker(500 * time.Millisecond)
	ticker := time.NewTicker(1 * time.Second)

	// go func() {
	for range ticker.C {
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(message),
			})

		if err != nil {
			log.Fatalf("%s: %s", "Failed to publish a message", err)
		}

		fmt.Println("publish success!")
	}
	// }()

	// time.Sleep(5 * time.Second)
	// ticker.Stop()
}

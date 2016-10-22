package main

import (
	"fmt"

	"github.com/op/go-logging"
	"github.com/streadway/amqp"
)

var log = logging.MustGetLogger("logger")

func Receive() {
	conn, err := amqp.Dial("amqp://guest:guest@cerbercam.cloudapp.net:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		"picturesQueue", // name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			fmt.Print("Received a message: %s", d.Body)
		}
	}()

	defer ch.Close()
	defer conn.Close()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Criticalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

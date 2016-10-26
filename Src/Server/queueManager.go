package main

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

// Receive data from queue
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

	i := 0
	for d := range msgs {
		i++
		msg := Deserialize(d.Body)
		log.Infof("Received a message: %s", *msg.Email)

		// open output file
		fo, err := os.Create(fmt.Sprintf("photo_%d.jpg", i))
		failOnError(err, "Failed to create file")

		_, err = fo.Write(msg.Photo)
		failOnError(err, "Failed to write to file")

		// close fo on exit and check for its returned error
		defer func() {
			err := fo.Close()
			failOnError(err, "Failed to close file")
		}()
	}

	defer ch.Close()
	defer conn.Close()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Criticalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

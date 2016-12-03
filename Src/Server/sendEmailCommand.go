package main

import (
	"fmt"

	"github.com/golang/protobuf/proto"
)

// HandleSendEmail - handles getting data from queue and sending them to the user
func HandleSendEmail() {
	log.Info("Checking for new data in emails queue...")

	emailManager := NewEmailManager()
	queue := queueManager{}
	msgs := queue.Receive(GlobalConfig.Queue.Responses)
	for d := range msgs {

		msg := &Response{}
		err := proto.Unmarshal(d.Body, msg)
		failOnError(err, "Cannot deserialize response")

		content := fmt.Sprintf("Cerber believes that your picture shows %s (probability %f%%)", *msg.Label, *msg.Probability)
		emailManager.Send(*msg.Email, "Reconginion results", content)
	}
}

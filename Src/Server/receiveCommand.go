package main

import (
	"github.com/golang/protobuf/proto"
)

// HandleReceiveCommand - handles getting data from queue and inserting to database
func HandleReceiveCommand() {
	log.Info("Checking for new data in queue...")

	queue := queueManager{}
	msgs := queue.Receive(GlobalConfig.Queue.Requests)
	for d := range msgs {

		msg := &Message{}
		err := proto.Unmarshal(d.Body, msg)
		failOnError(err, "Cannot deserialize request message")
		label, probability := Recognize(msg)
		log.Infof("Tensorflow results: label - {s} (%d)", label, probability)

		response := Response{msg.Email, &label, &probability, nil}
		responsebytes, err := proto.Marshal(&response)
		failOnError(err, "Cannot serialize response message")
		queue.Send(GlobalConfig.Queue.Responses, &responsebytes)
	}
}

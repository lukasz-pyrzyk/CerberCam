package main

// HandleReceiveCommand - handles getting data from queue and inserting to database
func HandleReceiveCommand() {
	log.Info("Checking for new data in queue...")

	serializer := serializer{}
	queue := queueManager{}
	msgs := queue.Receive(GlobalConfig.Queue.Requests)
	i := 0
	for d := range msgs {
		i++
		log.Infof("Processing message %d", i)
		msg := serializer.Deserialize(d.Body)
		label, probability := Recognize(msg)
		log.Infof("Tensorflow results: label - {s} (%d)", label, probability)

		response := Response{msg.Email, &label, &probability, nil}
		responsebytes := serializer.Serialize(&response)
		queue.Send(GlobalConfig.Queue.Responses, &responsebytes)
	}
}

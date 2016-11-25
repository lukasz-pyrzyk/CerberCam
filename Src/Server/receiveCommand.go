package main

// HandleReceiveCommand - handles getting data from queue and inserting to database
func HandleReceiveCommand() {
	log.Info("Checking for new data in queue...")

	serializer := serializer{}
	queue := queueManager{}
	msgs := queue.Receive(GlobalConfig.Queue.Topic)
	for d := range msgs {
		msg := serializer.Deserialize(d.Body)
		log.Infof("Received a message: %s", *msg.Email)

		InsertToDatabase(msg)
	}
}

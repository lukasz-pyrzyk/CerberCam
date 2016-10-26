package main

// HandleReceiveCommand - handles getting data from queue and inserting to database
func HandleReceiveCommand() {
	log.Info("Checking for new data in queue...")
	msgs := Receive("picturesQueue")
	for d := range msgs {
		msg := Deserialize(d.Body)
		log.Infof("Received a message: %s", *msg.Email)

		InsertToDatabase(msg)
	}
}

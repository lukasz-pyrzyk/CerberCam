package main

// HandleSendCommand - handles getting data from database and recognition with tensorflow
func HandleSendCommand() {
	log.Info("Checking for new data in database...")
	msgs := ReceiveFromDatabase()
	for _, msg := range msgs {
		log.Infof("Received a message: %s", *msg.Email)
		recognize(msg, "cerbercam.cloudapp.net:8888")
	}
}

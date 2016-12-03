package main

// HandleSendEmail - handles getting data from queue and sending them to the user
func HandleSendEmail() {
	log.Info("Checking for new data in emails queue...")

	email := NewEmailManager()
	email.send("lukasz.pyrzyk@gmail.com", "topic!", "message")
}

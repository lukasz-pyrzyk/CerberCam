package main

import "net/smtp"
import "strconv"

// HandleSendEmail - handles getting data from queue and sending them to the user
func HandleSendEmail() {
	log.Info("Checking for new data in emails queue...")
	// Set up authentication information.
	auth := smtp.PlainAuth("", GlobalConfig.Email.Login, GlobalConfig.Email.Password, GlobalConfig.Email.Host)

	to := []string{"temp@gmail.com"}
	msg := []byte("To: recipient@example.net\r\n" +
		"Subject: discount Gophers!\r\n" +
		"\r\n" +
		"This is the email body.\r\n")
	hostString := GlobalConfig.Email.Host + ":" + strconv.Itoa(GlobalConfig.Email.Port)
	err := smtp.SendMail(hostString, auth, GlobalConfig.Email.Login, to, msg)
	if err != nil {
		log.Fatal(err)
	}
}

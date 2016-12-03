package main

import (
	"fmt"
	"net/smtp"
	"strconv"
)

type emailManager struct {
}

func (email emailManager) send(to string, subject string, content string) {
	auth := smtp.PlainAuth("", GlobalConfig.Email.Login, GlobalConfig.Email.Password, GlobalConfig.Email.Host)
	recipients := []string{to}

	message := fmt.Sprintf("To: %s \r\n"+
		"Subject: %s !\r\n"+
		"\r\n"+
		"%s \r\n", to, subject, content)

	msg := []byte(message)
	hostString := GlobalConfig.Email.Host + ":" + strconv.Itoa(GlobalConfig.Email.Port)
	err := smtp.SendMail(hostString, auth, GlobalConfig.Email.Login, recipients, msg)
	failOnError(err, "Cannot send email")
}

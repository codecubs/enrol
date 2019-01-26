package email

import (
	"fmt"
	"net/mail"
)

type AdminEmail struct {
	Addr    mail.Address
	Subject string
	Body    string
}

func (e AdminEmail) Message() string {
	// from := mail.Address{"", "james@codecubs.co.uk"}
	// to := mail.Address{"", "james@codecubs.co.uk"}
	// subj := "New Enrol"
	// body := "email: " + event.Email

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = e.Addr.String()
	headers["To"] = e.Addr.String()
	headers["Subject"] = e.Subject

	// Setup message
	var message string
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + e.Body

	return message
}

func (e AdminEmail) FromAddr() string {
	return e.Addr.Address
}

func (e AdminEmail) ToAddr() string {
	return e.Addr.Address
}

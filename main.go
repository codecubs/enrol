package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Email string `json:"email"`
}

type MyResponse struct {
	Message string `json:"message:"`
	Ok      bool   `json:"ok"`
}

func HandleLambdaEvent(event MyEvent) (MyResponse, error) {
	from := mail.Address{"", "james@codecubs.co.uk"}
	to := mail.Address{"", "james@codecubs.co.uk"}
	subj := "New Enrol"
	body := "email: " + event.Email

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server
	servername := "smtp.mail.eu-west-1.awsapps.com:465"

	host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", "james@codecubs.co.uk", os.Getenv("password"), host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		return ReportIssue(event, err)
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return ReportIssue(event, err)
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		return ReportIssue(event, err)
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		return ReportIssue(event, err)
	}

	if err = c.Rcpt(to.Address); err != nil {
		return ReportIssue(event, err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return ReportIssue(event, err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return ReportIssue(event, err)
	}

	err = w.Close()
	if err != nil {
		return ReportIssue(event, err)
	}

	c.Quit()

	return MyResponse{
		Message: fmt.Sprintf("Sent enrollment email regarding: %s!", event.Email),
		Ok:      true,
	}, nil
}

func ReportIssue(event MyEvent, err error) (MyResponse, error) {
	fmt.Println(err.Error())
	return MyResponse{
		Message: fmt.Sprintf("Failed to send enrollment email regarding %s. Error: %s", event.Email, err.Error()),
		Ok:      false,
	}, err
}

func main() {
	lambda.Start(HandleLambdaEvent)
}

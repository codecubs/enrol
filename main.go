package main

import (
	"fmt"
	"net/mail"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/codecubs/enrol/email"
)

type MyEvent struct {
	Email string `json:"email"`
}

type MyResponse struct {
	Message string `json:"message:"`
	Ok      bool   `json:"ok"`
}

func HandleEnrolEvent(event MyEvent) (MyResponse, error) {
	m := email.AdminEmail{
		mail.Address{"", "james@codecubs.co.uk"},
		"New Enrol",
		"email: " + event.Email,
	}

	if err := email.Send(m); err != nil {
		return MyResponse{
			Message: fmt.Sprintf("Something went wrong: ", err.Error()),
			Ok:      false,
		}, err
	}

	return MyResponse{
		Message: fmt.Sprintf("Sent enrollment email regarding: %s!", event.Email),
		Ok:      true,
	}, nil
}

func main() {
	lambda.Start(HandleEnrolEvent)
}

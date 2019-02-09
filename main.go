package main

import (
	"fmt"
	"net/mail"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/codecubs/enrol/email"
)

type MyEvent struct {
	Parent  string `json:"parent"`
	Email   string `json:"email"`
	Student string `json:"student"`
	Group   string `json:"group"`
}

type MyResponse struct {
	Message string `json:"message:"`
	Ok      bool   `json:"ok"`
}

func HandleEnrolEvent(event MyEvent) (MyResponse, error) {

	contactEmail, err := email.SanitizeEmail(event.Email)
	if err != nil {
		return MyResponse{
			Message: fmt.Sprintf("Something is wrong with the email: %s", err.Error()),
			Ok:      false,
		}, err
	}

	m := email.AdminEmail{
		mail.Address{"", os.Getenv("adminemail")},
		"New Enrol",
		fmt.Sprintf("parent: %s, email:%s, student: %s, group %s", event.Parent, contactEmail, event.Student, event.Group),
	}

	if err := email.Send(m); err != nil {
		return MyResponse{
			Message: fmt.Sprintf("Something went wrong in enrolling: %s", err.Error()),
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

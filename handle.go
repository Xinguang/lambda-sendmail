package main

import (
	"net/mail"

	"github.com/Xinguang/lambda-sendmail/sendmail"
	"github.com/aws/aws-lambda-go/events"
	log "github.com/sirupsen/logrus"
	"github.com/xinguang/go-recaptcha"
)

type Contact struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

var (
	captcha *recaptcha.ReCAPTCHA
)
func init() {
	_captcha, err := recaptcha.New()
	if err != nil {
		log.Fatal("ReCAPTCHA:", err)
	}
	captcha = _captcha
}

// handle for verify a user's response to a reCAPTCHA challenge
func handle(contact Contact) (events.APIGatewayProxyResponse, error) {
	log.Debugf("Contact: %s", contact)
	err := captcha.Verify(contact.Token)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 400,
		}, nil
	}

	reply := mail.Address{Name: contact.Name, Address: contact.Email}
	_, err = sendmail.Send(reply, contact.Message)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "there is some errors",
			StatusCode: 400,
		}, err
	}
	return events.APIGatewayProxyResponse{
		Body:       contact.Message,
		StatusCode: 200,
	}, nil
}

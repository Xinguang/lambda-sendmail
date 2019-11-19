package recaptcha

import (
	"net/mail"

	"github.com/Xinguang/lambda-sendmail/sendmail"
	"github.com/aws/aws-lambda-go/events"
	log "github.com/sirupsen/logrus"
)

type Contact struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

// Handle for verify a user's response to a reCAPTCHA challenge
func Handle(contact Contact) (events.APIGatewayProxyResponse, error) {
	if !verify(contact.Token) {
		return events.APIGatewayProxyResponse{
			Body:       "timeout-or-duplicate",
			StatusCode: 400,
		}, nil
	}

	log.Info(contact)
	reply := mail.Address{Name: contact.Name, Address: contact.Email}
	message := sendmail.NewMessage()

	message.Add("Reply-To", reply.String())
	message.SetBody(contact.Message)

	log.Info("send email")
	err := sendmail.Send(*message)

	if err != nil {
		log.Fatal(err)
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

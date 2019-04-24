package recaptcha

import (
	"errors"
	"net/mail"

	"github.com/Xinguang/lambda-sendmail/sendmail"
	log "github.com/sirupsen/logrus"
)

type Contact struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

// Handle for verify a user's response to a reCAPTCHA challenge
func Handle(contact Contact) (string, error) {
	if !verify(contact.Token) {
		return "", errors.New("timeout-or-duplicate")
	}

	reply := mail.Address{Name: contact.Name, Address: contact.Email}
	message := sendmail.NewMessage()

	message.Add("Reply-To", reply.String())
	message.SetBody(contact.Message)

	log.Info("send email")
	err := sendmail.Send(*message)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return "verify", nil
}

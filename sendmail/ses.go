package sendmail

import (
	log "github.com/sirupsen/logrus"
	"net/mail"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type sesSender struct {
	client      *ses.SES
	destination *ses.Destination
	source      *string
	isHtml      bool
	subject     string
}

func newSESSender(sendmail *SendMail) *sesSender {
	// SES
	config := aws.NewConfig().WithRegion(sendmail.env.SesRegion)

	sender := &sesSender{}
	sender.client = ses.New(session.New(), config)

	sender.isHtml = sendmail.env.IsHtml
	sender.subject = sendmail.env.MailSubject
	sender.source = aws.String(sendmail.from.String())
	// ses.Destination
	sender.destination = &ses.Destination{
		// The recipients to place on the To: line of the message.
		ToAddresses: sender.parseAddressList(sendmail.toList),
		// The recipients to place on the CC: line of the message.
		CcAddresses: sender.parseAddressList(sendmail.ccList),
		// The recipients to place on the BCC: line of the message.
		BccAddresses: sender.parseAddressList(sendmail.bccList),
	}
	return sender
}
func (sender *sesSender) parseAddressList(emails []*mail.Address) []*string {
	if emails == nil || len(emails) == 0 {
		return nil
	}
	addressList := []*string{}
	for _, v := range emails {
		addressList = append(addressList, aws.String(v.String()))
	}
	return addressList
}

func (sender *sesSender) getSendEmailInput(reply mail.Address, message string) *ses.SendEmailInput {
	content := &ses.Content{
		Charset: aws.String("UTF-8"),
		Data:    aws.String(message),
	}
	body := &ses.Body{}
	if sender.isHtml {
		body.Html = content
	} else {
		body.Text = content
	}
	return &ses.SendEmailInput{
		Destination: sender.destination,
		Message: &ses.Message{
			Body: body,
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(sender.subject),
			},
		},
		ReplyToAddresses: []*string{
			aws.String(reply.String()),
		},
		Source: sender.source,
	}
}
func (sender *sesSender) Send(reply mail.Address, message string) (*string, error) {
	log.Info("reply", reply)
	log.Info("sender", sender)
	input := sender.getSendEmailInput(reply, message)
	result, err := sender.client.SendEmail(input)
	if err != nil {
		return nil, err
	}
	return result.MessageId, nil
}

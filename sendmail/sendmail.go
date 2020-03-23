package sendmail

import (
	log "github.com/sirupsen/logrus"
	"net/mail"
)

type SendMail struct {
	env ENV
	from    *mail.Address
	ccList  []*mail.Address
	bccList []*mail.Address
	toList  []*mail.Address
}

var (
	sendmail  = &SendMail{}
)

func init() {
	sendmail.env = envHolder
	var err error
	sendmail.from, err = mail.ParseAddress(envHolder.MailFrom)
	if err != nil {
		log.Fatal("MAIL_FROM:", err)
	}
	sendmail.toList, err = parseAddressList(envHolder.ToAddresses)
	if err != nil {
		log.Fatalf("TO_ADDRESSES: %s error: %s", envHolder.ToAddresses, err)
	}
	sendmail.ccList, err = parseAddressList(envHolder.CcAddresses)
	if err != nil {
		log.Fatal("CC_ADDRESSES:", err)
	}
	sendmail.bccList, err = parseAddressList(envHolder.BccAddresses)
	if err != nil {
		log.Fatal("BCC_ADDRESSES:", err)
	}
}

func parseAddressList(list string) ([]*mail.Address, error) {
	if len(list) == 0 {
		return nil, nil
	}
	return mail.ParseAddressList(list)
}

func Send(reply mail.Address, message string) (*string, error) {
	log.Info("sendmail", sendmail)

	log.Info("reply", reply)
	if len(sendmail.env.SmtpHost) == 0 {
		sender := newSESSender(sendmail)
		return sender.Send(reply, message)
	} else {
		sender := newSMTPSender(sendmail)
		err := sender.Send(reply, message)
		return nil, err
	}
	return nil, nil

	// if len(SMTP_HOST)>0 {
	// 	msg := sendmail.NewMessage()

	// 	msg.Add("Reply-To", reply.String())
	// 	msg.SetBody(message)

	// 	log.Info("send email")
	// 	err = sendmail.Send(*message)
	// }
}

package sendmail

import (
	log "github.com/sirupsen/logrus"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"
	"encoding/base64"
)

type smtpSender struct {
	header string
	boundary string
	body   string
	isHtml      bool

	mailAuth smtp.Auth 
	// The servername must include a port, as in "mail.example.com:smtp".
	servername string
	from string
	toList []string
}


func (sender *smtpSender) addHeader(field, value string) {
	if field == "boundary" {
		sender.header += fmt.Sprintf("%s\r\n", value)
	} else {
		sender.header += fmt.Sprintf("%s: %s\r\n", field, value)
	}
}

func (sender *smtpSender) setBody(body string) {
	contentType := "plain"
	if sender.isHtml {
		contentType = "html"
	}
	sender.body = fmt.Sprintf("Content-Type: text/%s; charset=\"utf-8\"\r\n", contentType)
	sender.body += "Content-Transfer-Encoding: base64\r\n"
	sender.body += base64.StdEncoding.EncodeToString([]byte(body))
}
// The msg parameter should be an RFC 822-style email with headers
// first, a blank line, and then the message body. The lines of msg
// should be CRLF terminated. The msg headers should usually include
// fields such as "From", "To", "Subject", and "Cc".  Sending "Bcc"
// messages is accomplished by including an email address in the to
// parameter but not including it in the msg headers.

func (sender *smtpSender) getContent() []byte {
	content := sender.header
	// // mixed
	// content += fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n", sender.boundary)
	// // split line
	// content += fmt.Sprintf("\r\n--%s\r\n", sender.boundary)
	// mail body
	content += sender.body
	
	return []byte(content)
}

// Send connects to the server at addr, switches to TLS if
// possible, authenticates with the optional mechanism a if possible,
// and then sends an email from address from, to addresses to, with
// message msg.
//
// The addresses in the to parameter are the SMTP RCPT addresses.
func (sender *smtpSender) Send(reply mail.Address, message string) error {
	log.Infof("Send with tls %s", reply)
	err := sender.sendWithTls(reply, message)
	if(err == nil) {
		return nil
	}
	log.Infof("Send %s", reply)
	sender.addHeader("Reply-To", reply.String())
	sender.setBody(message)
	return smtp.SendMail(sender.servername, sender.mailAuth, sender.from, sender.toList, sender.getContent())
}

func (sender *smtpSender) parseAddressList(list []*mail.Address) string {
	strList := []string{}
	for _, m := range list {
		strList = append(strList, m.Address)
	}
	return strings.Join(strList, ",")
}

func newSMTPSender(sendmail *SendMail) *smtpSender {
	sender := &smtpSender{}
	sender.isHtml = sendmail.env.IsHtml
	sender.addHeader("Subject", sendmail.env.MailSubject)
	sender.addHeader("From", sendmail.from.String())
	sender.addHeader("To", sender.parseAddressList(sendmail.toList))
	sender.addHeader("Cc", sender.parseAddressList(sendmail.ccList))
	sender.addHeader("Bcc", sender.parseAddressList(sendmail.bccList))
	sender.addHeader("MIME-Version", "1.0")

	sender.mailAuth = smtp.PlainAuth("", sendmail.env.SmtpUser, sendmail.env.SmtpPassword, sendmail.env.SmtpHost)
	// The servername must include a port, as in "mail.example.com:smtp".
	sender.servername = fmt.Sprintf("%s:%d", sendmail.env.SmtpHost, sendmail.env.SmtpPort)

	sender.from = sendmail.from.Address
	sender.toList = []string{}
	for _, to := range sendmail.toList {
		sender.toList = append(sender.toList, to.Address)
	}
	return sender
}

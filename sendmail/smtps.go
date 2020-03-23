package sendmail

import (
	log "github.com/sirupsen/logrus"
	"net/mail"
	"net/smtp"
	"net"
	"crypto/tls"
	"io"
)

func (sender *smtpSender) dial() (*smtp.Client, error) {
    host, _, _ := net.SplitHostPort(sender.servername)
	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}
	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", sender.servername, tlsconfig)
	if err != nil {
		return nil, err
	}
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// SendMail connects to the server at addr, switches to TLS if
// possible, authenticates with the optional mechanism a if possible,
// and then sends an email from address from, to addresses to, with
// message msg.
//
// The addresses in the to parameter are the SMTP RCPT addresses.
func (sender *smtpSender) sendWithTls(reply mail.Address, message string) error {
	client, err := sender.dial()
	if err != nil {
		return err
	}
	defer client.Close()

	w := sender.writeCloser(client, sender.from)
	_, err = w.Write(sender.getContent())
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return client.Quit()
}

func (sender *smtpSender) writeCloser(client *smtp.Client, from string) io.WriteCloser {
	// Auth
	if err := client.Auth(sender.mailAuth); err != nil {
		log.Fatal(err)
	}
	// To && From
	if err := client.Mail(from); err != nil {
		log.Fatal(err)
	}
	for _, to := range sender.toList {
		if err := client.Rcpt(to); err != nil {
			log.Fatal(err)
		}
	}
	// Data
	w, err := client.Data()
	if err != nil {
		log.Fatal(err)
	}
	return w
}
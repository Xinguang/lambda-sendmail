package sendmail

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/mail"
	"net/smtp"
	"os"

	log "github.com/sirupsen/logrus"
)

var (
	SENDER_NAME     = os.Getenv("SENDER_NAME")
	SENDER_USER     = os.Getenv("SENDER_USER")
	SENDER_PASSWORD = os.Getenv("SENDER_PASSWORD")
	SENDER_HOST     = os.Getenv("SENDER_HOST")
	SENDER_PORT     = os.Getenv("SENDER_PORT") // default 587
	RECIPIENT       = os.Getenv("RECIPIENT")   // name<mail@server.addr>
)

// The recipient address
var recipient = func() *mail.Address {
	address, err := mail.ParseAddress(RECIPIENT)
	if err != nil {
		log.Fatal(err)
	}
	return address
}()

func Send(msg Message) error {
	// The servername must include a port, as in "mail.example.com:smtp".
	servername := fmt.Sprintf("%s:%s", SENDER_HOST, SENDER_PORT)
	auth := smtp.PlainAuth("", SENDER_USER, SENDER_PASSWORD, SENDER_HOST)
	recipients := []string{recipient.Address}
	msg.Add("To", recipient.String())
	from := mail.Address{Name: SENDER_NAME, Address: SENDER_USER}
	msg.Add("From", from.String())

	log.Infof("recipients %s", recipient.String())

	return smtp.SendMail(servername, auth, from.Address, recipients, msg.getContent())
}

// SendMail connects to the server at addr, switches to TLS if
// possible, authenticates with the optional mechanism a if possible,
// and then sends an email from address from, to addresses to, with
// message msg.
//
// The addresses in the to parameter are the SMTP RCPT addresses.
func SendMail(from string, msg Message) error {
	client := dial()
	defer client.Close()

	w := writeCloser(client, from)
	_, err := w.Write(msg.getContent())
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return client.Quit()
}

func dial() *smtp.Client {
	// The servername must include a port, as in "mail.example.com:smtp".
	servername := fmt.Sprintf("%s:%s", SENDER_HOST, SENDER_PORT)
	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         SENDER_HOST,
	}
	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		log.Fatal(err)
	}
	client, err := smtp.NewClient(conn, SENDER_HOST)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

var writeCloser = func(client *smtp.Client, from string) io.WriteCloser {
	auth := smtp.PlainAuth("", SENDER_USER, SENDER_PASSWORD, SENDER_HOST)
	// Auth
	if err := client.Auth(auth); err != nil {
		log.Fatal(err)
	}
	// To && From
	if err := client.Mail(from); err != nil {
		log.Fatal(err)
	}
	if err := client.Rcpt(recipient.Address); err != nil {
		log.Fatal(err)
	}
	// Data
	w, err := client.Data()
	if err != nil {
		log.Fatal(err)
	}
	return w
}

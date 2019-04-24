package sendmail

import (
	"encoding/base64"
	"fmt"
	"os"
)

var (
	SENDMAIL_SUBJECT  = os.Getenv("SENDMAIL_SUBJECT")
	SENDMAIL_MAILTYPE = os.Getenv("SENDMAIL_MAILTYPE") // html , plain ....
)

type Message struct {
	header map[string]string
	body   string
}

func NewMessage() *Message {
	message := &Message{}
	message.header = make(map[string]string)
	message.Add("Subject", SENDMAIL_SUBJECT)
	message.Add("MIME-Version", "1.0")
	if len(SENDMAIL_MAILTYPE) > 0 {
		message.Add("Content-Type", "text/"+SENDMAIL_MAILTYPE+"; charset=\"utf-8\"")
	} else {
		message.Add("Content-Type", "text/plain; charset=\"utf-8\"")
	}
	message.Add("Content-Transfer-Encoding", "base64")
	return message
}

func (msg *Message) Add(field, value string) {
	msg.header[field] = value
}
func (msg *Message) SetBody(body string) {
	msg.body = body
}

// The msg parameter should be an RFC 822-style email with headers
// first, a blank line, and then the message body. The lines of msg
// should be CRLF terminated. The msg headers should usually include
// fields such as "From", "To", "Subject", and "Cc".  Sending "Bcc"
// messages is accomplished by including an email address in the to
// parameter but not including it in the msg headers.

func (msg *Message) getContent() []byte {
	content := ""
	for k, v := range msg.header {
		content += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	content += "\r\n" + base64.StdEncoding.EncodeToString([]byte(msg.body))
	return []byte(content)
}

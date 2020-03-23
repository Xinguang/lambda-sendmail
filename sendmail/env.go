package sendmail

import (
	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"
)

type ENV struct {
	MailSubject  string `env:"MAIL_SUBJECT" envDefault:"Test email"`
	SesRegion    string `env:"SES_REGION" envDefault:"us-west-1"`
	IsHtml       bool   `env:"IS_HTML" envDefault:true`
	MailFrom     string `env:"MAIL_FROM"`
	CcAddresses  string `env:"CC_ADDRESSES"`
	BccAddresses string `env:"BCC_ADDRESSES"`
	ToAddresses  string `env:"TO_ADDRESSES"`
	SmtpUser     string `env:"SMTP_USER"`
	SmtpPassword string `env:"SMTP_PASSWORD"`
	SmtpHost     string `env:"SMTP_HOST"`
	SmtpPort     int    `env:"SMTP_PORT" envDefault:587`
}
var (
	envHolder = ENV{}
)

func init() {
	err := env.Parse(&envHolder)
	if err != nil {
		log.Fatal(err)
	}
}
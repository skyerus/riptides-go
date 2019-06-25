package email

import (
	"github.com/mailgun/mailgun-go"
	"os"
)

func NewMailgun() *mailgun.MailgunImpl {
	mg := mailgun.NewMailgun(os.Getenv("MAILGUN_DOMAIN"), os.Getenv("MAILGUN_API_KEY"))
	mg.SetAPIBase("https://api.eu.mailgun.net/v3")

	return mg
}

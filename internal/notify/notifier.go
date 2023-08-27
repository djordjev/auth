package notify

import (
	"fmt"

	"github.com/djordjev/auth/internal/utils"
	"github.com/mailjet/mailjet-apiv3-go/v4"
)

type MailjetNotifier struct {
	client *mailjet.Client
	config utils.Config
}

func NewMailjetNotifier(config utils.Config) MailjetNotifier {
	client := mailjet.NewMailjetClient(config.Mailjet.ApiKey, config.Mailjet.SecretKey)
	notifier := MailjetNotifier{client: client, config: config}

	return notifier
}

func (mj MailjetNotifier) Send(to string, subject string, text string, html string) error {
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From:     &mailjet.RecipientV31{Email: mj.config.Sender, Name: "Authentication"},
			To:       &mailjet.RecipientsV31{mailjet.RecipientV31{Email: to}},
			Subject:  subject,
			TextPart: text,
			HTMLPart: html,
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}
	_, err := mj.client.SendMailV31(&messages)
	if err != nil {
		return fmt.Errorf("unable to send email message to %s", to)
	}

	return nil
}

type SilentNotifier struct{}

func (dn SilentNotifier) Send(to string, subject string, text string, html string) error {
	fmt.Printf("Sending email to %s with subject %s and body %s\n", to, subject, text)

	return nil
}

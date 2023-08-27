package domain

type Notifier interface {
	Send(to string, subject string, text string, html string) error
}

var verifyTextTemplate = "Hello %s, thanks for creating new account. Please follow this link %s to verify it."
var verifyHtmlTemplate = `<h2>Hello %s</h2><p>thanks for creating a new account. Please click the button below to verify it</p><a href="%s">Verify</a>`

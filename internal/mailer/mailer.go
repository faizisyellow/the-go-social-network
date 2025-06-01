package mailer

import "embed"

const (
	// the name of the sender email
	FromName            = "The Go Social Network"
	maxRetries          = 3
	UserWelcomeTemplate = "user_invitation.tmpl"
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(TemplateFile, username, mail string, data any, isSandbox bool) (int, error)
}

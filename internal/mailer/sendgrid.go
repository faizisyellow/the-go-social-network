package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendgrid(apikey, fromEmail string) *SendGridMailer {
	client := sendgrid.NewSendClient(apikey)

	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apikey,
		client:    client,
	}
}

func (s *SendGridMailer) Send(TemplateFile, username, email string, data any, isSandbox bool) (int, error) {
	from := mail.NewEmail(FromName, s.fromEmail)

	to := mail.NewEmail(username, email)

	// templating and building
	tmpl, err := template.ParseFS(FS, "templates/"+TemplateFile)
	if err != nil {
		return -1, err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return -1, err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return -1, err
	}

	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())

	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			// if sandbox enable it will response(200) mean the email send is success but wont send it to the recipient
			// if not enable it will response(202 it means the email success and send it to recipient)
			Enable: &isSandbox,
		},
	})

	var ErrRetry error
	for i := 0; i < maxRetries; i++ {
		response, ErrRetry := s.client.Send(message)
		if ErrRetry != nil {

			// exponential backoff
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		return response.StatusCode, nil
	}

	// if somehow error, or the attemps is send over max retries
	return -1, fmt.Errorf("failed to send email after %d attemps error: %v ", maxRetries, ErrRetry)
}

package channel

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"time"

	"github.com/lz1marine/notification-service/pkg/queue"

	gomail "gopkg.in/mail.v2"
)

// TemplatePreview is used to preview email templates
type TemplatePreview struct {
	Name, Email, Message string
}

type email struct {
	sender                    string
	shortTimeout, longTimeout time.Duration

	dialer *gomail.Dialer
}

// NewEmailChannel creates a new email channel
func NewEmailChannel(host string, port int, username string, password string) *email {
	e := &email{
		sender:       username,
		shortTimeout: 5 * time.Second,
		longTimeout:  60 * time.Second,
		dialer: gomail.NewDialer(
			host,
			port,
			username,
			password,
		),
	}

	e.dialer.Timeout = e.shortTimeout
	return e
}

// Name returns the name of the channel
func (e *email) Name() string {
	return "email"
}

// Notify sends an email
func (e *email) Notify(m *queue.Message) error {
	message, err := e.prepare(m.Recepients, m.Message, m.Subject, m.Template)
	if err != nil {
		return err
	}

	// TODO: log info
	fmt.Printf("message: %+v\n", message)

	err = e.send(message, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (e *email) prepare(to []string, message string, subject *string, t *template.Template) (*gomail.Message, error) {
	m := gomail.NewMessage()
	m.SetHeader("From", e.sender)
	m.SetHeader("To", to...)
	if subject != nil {
		m.SetHeader("Subject", *subject)
	}

	m = e.attachBody(m, message, to, t)
	return m, nil
}

func (e *email) attachBody(m *gomail.Message, message string, to []string, t *template.Template) *gomail.Message {
	if t != nil {
		toEmail := to[0]
		if len(to) > 1 {
			toEmail = "Everyone"
		}

		preview := TemplatePreview{
			Name:    "temp_placeholder", // TODO: get name
			Email:   toEmail,
			Message: message,
		}

		var tpl bytes.Buffer
		if err := t.Execute(&tpl, preview); err != nil {
			fmt.Println(err)
		}

		m.SetBody("text/html", tpl.String())
	} else {
		m.SetBody("text/plain", message)
	}

	return m
}

func (e *email) send(m *gomail.Message, start time.Time) error {
	now := time.Now()
	if now.Sub(start) >= e.longTimeout {
		return errors.New("failed to send email")
	}

	// TODO: potential hang here. We should add a context with a timeout that will kill the goroutine
	err := e.dialer.DialAndSend(m)
	if err != nil {
		fmt.Printf("failed to send email: %v\n", err)
		time.Sleep(e.shortTimeout * time.Second)
		return e.send(m, start)
	}

	// TODO: log info
	fmt.Println("email sent")
	return nil
}

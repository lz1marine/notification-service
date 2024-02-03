package channels

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/mail"
	"sync"
	"time"

	"github.com/lz1marine/notification-service/pkg/queue"

	gomail "gopkg.in/mail.v2"
)

type TemplatePreview struct {
	Name, Email, Message string
}

type email struct {
	sender                    string
	maxConnections            uint
	shortTimeout, longTimeout time.Duration

	dialer *gomail.Dialer
}

func NewEmailChannel(host string, port int, username string, password string, maxConnections uint) *email {
	e := &email{
		sender:         username,
		shortTimeout:   5 * time.Second,
		longTimeout:    60 * time.Second,
		maxConnections: maxConnections,
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

func (e *email) Name() string {
	return "email"
}

func (e *email) Notify(message *queue.Message) error {
	var wg sync.WaitGroup
	maxConCh := make(chan struct{}, e.maxConnections)

	errorsFound := false
	for _, em := range message.Recepients {
		if !isValidEmail(em) {
			continue
		}
		curEm := em
		message, err := e.prepare(curEm, message.Message, message.Title, message.Template)
		if err != nil {
			return err
		}

		// TODO: log info
		fmt.Printf("message: %+v\n", message)
		fmt.Printf("sending to: %s\n", curEm)

		// Start in async up to e.maxConnections
		wg.Add(1)
		maxConCh <- struct{}{}

		go func() {
			defer func() {
				<-maxConCh
			}()
			defer wg.Done()

			err := e.send(message, time.Now())
			if err != nil {
				errorsFound = true
				fmt.Printf("failed to send email to %s: %v\n", curEm, err)
			}
		}()

	}

	wg.Wait()

	if errorsFound {
		// TODO: change message state to failed, requeue
		return fmt.Errorf("failed to send emails")
	}

	// TODO: change message state to completed
	return nil
}

func (e *email) prepare(to, message string, subject *string, t *template.Template) (*gomail.Message, error) {
	m := gomail.NewMessage()
	m.SetHeader("From", e.sender)
	m.SetHeader("To", to)
	if subject != nil {
		m.SetHeader("Subject", *subject)
	}

	m = e.attachBody(m, to, message, t)
	return m, nil
}

func (e *email) attachBody(m *gomail.Message, to, message string, t *template.Template) *gomail.Message {
	if t != nil {
		preview := TemplatePreview{
			Name:    "temp_placeholder", // TODO: get name
			Email:   to,
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
	if time.Now().Sub(start) >= e.longTimeout {
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

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		fmt.Printf("failed to parse email %s: %v", email, err)
		return false
	}

	return true
}

package queue

import "html/template"

// TODO: move to another package

// Message is a message for a channel
type Message struct {
	// Message is the message body
	Message string

	// Recepients is the list of recepients
	Recepients []string

	// Subject is the message subject
	Subject *string

	// Template is the message template
	Template *template.Template
}

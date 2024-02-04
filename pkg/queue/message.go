package queue

import "html/template"

// TODO: move to another package

type Message struct {
	Message    string
	Recepients []string
	Subject    *string
	Template   *template.Template
}

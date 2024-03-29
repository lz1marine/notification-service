package client

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/go-redis/redis/v8"
	"github.com/lz1marine/notification-service/pkg/api"
)

// TemplateReader reads the templates from a source
type TemplateReader interface {
	// Read reads the template
	Read(templateID *string) (*template.Template, error)
}

// TemplateWriter writes the templates to a destination
type TemplateWriter interface {
	// Write writes the template
	Write(temaplte *api.Template) error
}

// TemplateReadWriter reads and writes the templates
type TemplateReadWriter interface {
	TemplateReader
	TemplateWriter
}

// RedisTemplate is used to read and write templates
type RedisTemplate struct {
	client *redis.Client
}

// NewRedisTemplate creates a new RedisTemplate
func NewRedisTemplate(endpoint, password string, db int) *RedisTemplate {
	options := &redis.Options{
		Addr:     endpoint,
		DB:       db,
		Password: password,
	}

	client := redis.NewClient(options)
	return &RedisTemplate{
		client: client,
	}
}

// Read reads the template
func (rt *RedisTemplate) Read(templateID *string) (*template.Template, error) {
	if templateID == nil {
		return nil, nil
	}

	val := rt.client.Get(context.Background(), *templateID)
	if val.Err() != nil {
		return nil, val.Err()
	}

	var t api.Template
	err := json.Unmarshal([]byte(val.Val()), &t)
	if err != nil {
		fmt.Printf("failed to unmarshal template: %v", err)
		return nil, err
	}

	res, err := template.New("template").Parse(t.Template)
	if err != nil {
		fmt.Printf("failed to parse template: %v", err)
		return nil, err
	}

	fmt.Printf("parsed template: %s\n", t.Template)
	return res, nil
}

// Write writes the template
func (rt *RedisTemplate) Write(temaplte *api.Template) error {
	// TODO: implement
	return nil
}

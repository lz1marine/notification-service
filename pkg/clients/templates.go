package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/go-redis/redis/v8"
	"github.com/lz1marine/notification-service/pkg/entities"
)

// TemplateReader reads the templates from a source
type TemplateReader interface {
	Read(templateID *string) (*template.Template, error)
}

// TemplateWriter writes the templates to a destination
type TemplateWriter interface {
	Write(temaplte entities.Templates) error
}

// TemplateReadWriter reads and writes the templates
type TemplateReadWriter interface {
	TemplateReader
	TemplateWriter
}

type RedisTemplate struct {
	client *redis.Client
}

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

func (rt *RedisTemplate) Read(templateID *string) (*template.Template, error) {
	if templateID == nil {
		return nil, nil
	}

	val := rt.client.Get(context.Background(), *templateID)
	if val.Err() != nil {
		return nil, val.Err()
	}

	t := entities.Templates{}
	json.Unmarshal([]byte(val.Val()), &t)

	res, err := template.New("template").Parse(t.Template)
	if err != nil {
		fmt.Printf("failed to parse template: %v", err)
	}

	fmt.Printf("parsed template: %s\n", t.Template)
	return res, nil
}

func (rt *RedisTemplate) Write(temaplte entities.Templates) error {
	// TODO: implement
	return nil
}

package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	apiv1 "github.com/lz1marine/notification-service/api/v1"
)

// BackupMessageSender sends the message to a destination
type BackupMessageSender interface {
	// Send sends the message
	Send(req *apiv1.NotificationRequest) error
}

// BackupMessageRemover removes the message
type BackupMessageRemover interface {
	// Remove removes the message
	Remove(req *apiv1.NotificationRequest) error
}

// BackupMessageReadWriter reads and writes the backup message.
type BackupMessageReadWriter interface {
	BackupMessageSender
	BackupMessageRemover
}

// RedisBackup is used to read and write backup messages
type RedisBackup struct {
	client *redis.Client
}

// NewRedisBackup creates a new RedisBackup
func NewRedisBackup(endpoint, password string, db int) *RedisBackup {
	options := &redis.Options{
		Addr:     endpoint,
		DB:       db,
		Password: password,
	}

	client := redis.NewClient(options)
	return &RedisBackup{
		client: client,
	}
}

// Send sends the message
func (rb *RedisBackup) Send(req *apiv1.NotificationRequest) error {
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	status := rb.client.Set(context.Background(), req.ID, string(reqBytes), 0)
	if status.Err() != nil {
		return status.Err()
	}

	return nil
}

// Remove removes the message
func (rb *RedisBackup) Remove(req *apiv1.NotificationRequest) error {
	deled := rb.client.Del(context.Background(), req.ID)
	if deled.Err() != nil {
		return deled.Err()
	}

	if deled.Val() != 1 {
		errStr := fmt.Sprintf("failed to delete message: %s, %d", req.ID, deled.Val())
		return errors.New(errStr)
	}

	return nil
}

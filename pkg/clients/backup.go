package clients

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	apiv1 "github.com/lz1marine/notification-service/api/v1"
)

// BackupMessageSender sends the message to a destination
type BackupMessageSender interface {
	Send(req *apiv1.NotificationRequest) error
}

// BackupMessageRemover removes the message
type BackupMessageRemover interface {
	Remove(req *apiv1.NotificationRequest) error
}

// BackupMessageReadWriter reads and writes the backup message.
type BackupMessageReadWriter interface {
	BackupMessageSender
	BackupMessageRemover
}

type RedisBackup struct {
	client *redis.Client
}

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

func (rb *RedisBackup) Send(req *apiv1.NotificationRequest) error {
	status := rb.client.Set(context.Background(), req.ID, req, 0)
	if status.Err() != nil {
		return status.Err()
	}

	return nil
}

func (rb *RedisBackup) Remove(req *apiv1.NotificationRequest) error {
	deled := rb.client.Del(context.Background(), req.ID)
	if deled.Err() != nil {
		return deled.Err()
	}

	if deled.Val() != 1 {
		return errors.New(fmt.Sprintf("failed to delete message: %s, %d", req.ID, deled.Val()))
	}

	return nil
}

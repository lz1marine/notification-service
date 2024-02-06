package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	apiv1 "github.com/lz1marine/notification-service/api/v1"
)

// Redis is used to push and pop messages
type Redis struct {
	retries int
	client  *redis.Client
}

// NewRedis creates a new Redis
func NewRedis(endpoint, password string, db, retries int) *Redis {
	options := &redis.Options{
		Addr:     endpoint,
		DB:       db,
		Password: password,
	}

	client := redis.NewClient(options)
	return &Redis{
		client: client,
	}
}

// Push adds a message to the queue
func (r *Redis) Push(req *apiv1.NotificationRequest, channel string) error {
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	return r.tryPush(string(reqBytes), channel, 0)
}

// Pop returns and removes the next message from the queue
func (r *Redis) Pop(channel string) (string, error) {
	popped := r.client.RPop(context.Background(), channel)
	res, err := popped.Result()
	if err != nil {
		return "", err
	}

	return res, nil
}

func (r *Redis) tryPush(req, channel string, try int) error {
	push := r.client.LPush(context.Background(), channel, req)
	if push.Err() != nil {
		if try < r.retries {
			fmt.Printf("failed to push message: %v\n", push.Err())
			time.Sleep(time.Second)
			return r.tryPush(req, channel, try+1)
		}

		return push.Err()
	}

	return nil
}

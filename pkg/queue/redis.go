package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	apiv1 "github.com/lz1marine/notification-service/api/v1"
)

type Redis struct {
	retries int
	client  *redis.Client
}

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

func (r *Redis) Push(req *apiv1.NotificationRequest, channel string) error {
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	return r.tryPush(string(reqBytes), channel, 0)
}

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
			time.Sleep(time.Second)
			return r.tryPush(req, channel, try+1)
		}

		return push.Err()
	}

	return nil
}

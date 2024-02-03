package queue

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
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

func (r *Redis) Push(req, channel string) error {
	return r.tryPush(req, channel, 0)
}

func (r *Redis) Pop(channel string) (string, error) {
	popped := r.client.BRPop(context.Background(), 30*time.Second, channel)
	res, err := popped.Result()
	if err != nil {
		return "", err
	}

	if len(res) != 1 {
		return "", errors.New("expected 1 popped element in redis")
	}

	return res[0], nil
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
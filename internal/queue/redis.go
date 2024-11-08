package queue

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var client redis.Client

func Init() error {
	ctx := context.Background()

	connection := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	_, err := connection.Ping(ctx).Result()
	if err != nil {
		return err
	}

	client = *connection

	return nil
}

package queue

import (
	"context"
	"fmt"
	"log"
	"time"

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

func InsertIntoQueue(queueName, item string) error {
	ctx := context.Background()

	err := client.LPush(ctx, queueName, item).Err()
	if err != nil {
		return fmt.Errorf("erro ao inserir item na fila: %w", err)
	}

	fmt.Printf("Item '%s' inserido na fila '%s'\n", item, queueName)

	return nil
}

func ConsumeQueue(queueName string) {
	ctx := context.Background()

	for {
		result, err := client.BLPop(ctx, 0*time.Second, queueName).Result()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			log.Printf("Erro ao consumir fila: %v", err)
			break
		}

		if len(result) > 1 {
			item := result[1]
			fmt.Printf("Processando item: %s\n", item)
		}
	}
}

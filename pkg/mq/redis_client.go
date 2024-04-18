package mq

import (
	"context"
	"inttest-runtime/pkg/utils"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisPubSub struct {
	client *redis.Client
}

func ConnectRedisPubSub(addr string, db int, password string) (*RedisPubSub, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       db,
		Password: password,
	})
	return &RedisPubSub{
		client: client,
	}, nil
}

func (pubSub RedisPubSub) Publish(ctx context.Context, topic string, message []byte) error {
	return pubSub.client.Publish(ctx, topic, message).Err()
}

func (pubSub RedisPubSub) Subscribe(ctx context.Context, topic string, consumer func(message []byte) error) {
	ps := pubSub.client.Subscribe(ctx, topic)
	// todo: unsubscribe
	go func() {
		for {
			anyResult, err := ps.Receive(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			switch m := anyResult.(type) {
			case *redis.Message:
				if err := consumer(utils.S2B(m.Payload)); err != nil {
					log.Printf("error while handling topic %s: %v", topic, err)
				}
			}
		}
	}()
}

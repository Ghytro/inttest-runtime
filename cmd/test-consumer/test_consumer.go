package main

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer client.Close()

	pubSub := client.Subscribe(context.Background(), "my_topic")

	for {
		anyResult, err := pubSub.Receive(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		switch m := anyResult.(type) {
		case *redis.Message:
			log.Printf("got a new message from channel: %s", m.Payload)
		}
	}
}

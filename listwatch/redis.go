package listwatch

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

func subscribe(topic string) *redis.PubSub {
	return client.Subscribe(ctx, topic)
}

// Publish is for testing now
func Publish(topic string, msg interface{}) {
	client.Publish(ctx, topic, msg)
}

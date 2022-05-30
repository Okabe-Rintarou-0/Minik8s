package listwatch

import (
	"context"
	"github.com/go-redis/redis/v8"
	"minik8s/global"
)

var ctx = context.Background()

var client = redis.NewClient(&redis.Options{
	Addr:     global.Host + ":6379",
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

package listwatch

import (
	"fmt"
	"github.com/go-redis/redis/v8"
)

type WatchHandler func(message *redis.Message)

// Watch infinitely watches a given topic, please call it
// by goroutine
func Watch(topic string, handler WatchHandler) {
	sub := subscribe(topic)
	fmt.Println("Subscribe", topic)
	for msg := range sub.Channel() {
		fmt.Printf("Received from %s: %s\n", msg.Channel, msg.Payload)
		handler(msg)
	}
}

package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var bgCtx = context.Background()
var cancel context.CancelFunc

func consumeHigherCpu(ctx context.Context) {
	for i := 0; i < 5; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}
			}
		}()
	}
}

func higher(c *gin.Context) {
	if cancel != nil {
		return
	}

	var ctx context.Context
	ctx, cancel = context.WithCancel(bgCtx)
	consumeHigherCpu(ctx)
	c.String(http.StatusOK, "higher cpu utilization!")
}

func lower(c *gin.Context) {
	if cancel != nil {
		cancel()
	}
	c.String(http.StatusOK, "lower cpu utilization!")
}

func main() {
	r := gin.Default()
	r.GET("/higher", higher)
	r.GET("/lower", lower)
	log.Fatal(r.Run(":8090"))
}

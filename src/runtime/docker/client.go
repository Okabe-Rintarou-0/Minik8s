package docker

import (
	"context"
	"github.com/docker/docker/client"
)

var Ctx = context.Background()

var Client = newClient()

func newClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	return cli
}

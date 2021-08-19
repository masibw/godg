package docker

import (
	"log"

	"github.com/docker/docker/client"
)

var (
	cli *client.Client
)

func init() {
	var err error
	cli, err = client.NewEnvClient()
	if err != nil {
		log.Fatalf("failed to get NewEnvClient: %v", err)
	}
}

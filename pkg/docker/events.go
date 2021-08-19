package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
)

// NewEventWatcher starts monitoring docker events.
func NewEventWatcher() (msg <-chan events.Message, err <-chan error) {
	filter := filters.NewArgs()
	filter.Add("type", "container")
	filter.Add("event", "start")

	msg, err = cli.Events(context.Background(), types.EventsOptions{Filters: filter})
	return
}

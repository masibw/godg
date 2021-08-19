package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/masibw/godg/pkg/docker"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "godg"
	app.Usage = "Graphical representation of docker container startup order"
	app.Flags = []cli.Flag{
		&cli.DurationFlag{
			Name:    "time",
			Aliases: []string{"t"},
			Value:   30 * time.Second,
			Usage:   "Waiting time for container starts",
		},
	}
	app.Action = func(c *cli.Context) error {
		log.Println("Starting godg")
		log.Println("Waiting for", c.Duration("time"))
		start := time.Now()
		go func() {
			// コンテナを監視する
			msgChan, errChan := docker.NewEventWatcher()
			for {
				select {
				case msg := <-msgChan:
					if msg.Action == "start" {
						fmt.Println("container name:", msg.Actor.Attributes["name"], "started in",time.Since(start))
					}
				case err := <-errChan:
					log.Fatalf("err: %v", err)
					return
				}
			}
		}()
		time.Sleep(c.Duration("time"))
		log.Printf("finished")
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}

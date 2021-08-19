package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/masibw/godg/pkg/docker"

	graphviz "github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
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
		log.Println("Monitoring for", c.Duration("time"))

		g := graphviz.New()
		graph, err := g.Graph()
		if err != nil {
			log.Fatalf("failed to graph: %v", err)
		}
		defer func() {
			if err := graph.Close(); err != nil {
				log.Fatalf("failed to close graph: %v", err)
			}
		}()

		start := time.Now()
		go func() {
			// コンテナを監視する
			msgChan, errChan := docker.NewEventWatcher()
			var beforeNode *cgraph.Node
			for {
				select {
				case msg := <-msgChan:
					if msg.Action == "start" {
						containerName := msg.Actor.Attributes["name"]
						fmt.Println("container name:", containerName, "started in", time.Since(start))
						var err error
						var node *cgraph.Node
						node, err = graph.CreateNode(containerName)
						if beforeNode != nil {
							_, err = graph.CreateEdge(containerName, beforeNode, node)
							if err != nil {
								log.Fatalf("failed to create node name: %s, err: %v", containerName, err)
							}
						}
						node.SetLabel(fmt.Sprintf("%s\n%s", containerName, time.Since(start)))
						beforeNode = node
					}
				case err := <-errChan:
					log.Fatalf("err: %v", err)
					return
				}
			}
		}()
		time.Sleep(c.Duration("time"))
		if err := g.RenderFilename(graph, graphviz.PNG, "./start-up.png"); err != nil {
			log.Fatalf("failed to output png image: %v", err)
		}
		log.Printf("finished")
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}

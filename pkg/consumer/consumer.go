package consumer

import (
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
)

type Consumer struct {
	Taskserver  *machinery.Server
	AppURL      string
	Concurrency int
}

func NewConsumer() (*Consumer, error) {
	c := &Consumer{
		AppURL:      "http://localhost:8081",
		Concurrency: 1,
	}

	taskserver, err := machinery.NewServer(&config.Config{
		Broker:        "redis://localhost:6379",
		ResultBackend: "redis://localhost:6379",
	})
	if err != nil {
		return nil, err
	}
	c.Taskserver = taskserver

	c.Taskserver.RegisterTasks(map[string]interface{}{
		"handyman": c.handleTask,
	})
	return c, nil
}

func (c *Consumer) Start() error {
	worker := c.Taskserver.NewWorker("handyman", c.Concurrency)
	err := worker.Launch()
	if err != nil {
		return err
	}

	return nil
}

func (c *Consumer) handleTask() error {
	// send task message in HTTP req to app server
	// get the response
	// do something with response

	return nil
}

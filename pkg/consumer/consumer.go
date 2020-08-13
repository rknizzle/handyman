package consumer

import (
	"bytes"
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"net/http"
)

type Consumer struct {
	Taskserver  *machinery.Server
	AppURL      string
	Concurrency int
	Client      HTTPClient
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewConsumer() (*Consumer, error) {
	c := &Consumer{
		AppURL:      "http://localhost:8081",
		Concurrency: 1,
		Client:      &http.Client{},
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

func (c *Consumer) handleTask(taskMessage string) (string, error) {
	// send task message in HTTP req to app server
	_, err := c.sendRequestToApp(taskMessage)
	if err != nil {
		return "", err
	}
	// get the response
	// do something with response

	return "returnVal", nil
}

func (c *Consumer) sendRequestToApp(taskMessage string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodPost, c.AppURL, bytes.NewReader([]byte(taskMessage)))
	if err != nil {
		return nil, err
	}
	return c.Client.Do(request)
}

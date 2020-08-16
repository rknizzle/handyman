package consumer

import (
	"bytes"
	"errors"
	"fmt"
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

// return an error when the HTTP request responds with a non 2xx status code
func (c *Consumer) handleTask(taskMessage string) (string, error) {
	// send task message in HTTP req to app server
	resp, err := c.sendRequestToApp(taskMessage)
	if err != nil {
		return "", err
	}

	// if non-2xx status code
	if string(resp.Status[0]) != "2" {
		return "", errors.New(fmt.Sprintf("Recieved non-successful status code %s", resp.Status))
	}

	// get JSON response
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	jsonBody := buf.String()

	return jsonBody, nil
}

func (c *Consumer) sendRequestToApp(taskMessage string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodPost, c.AppURL, bytes.NewReader([]byte(taskMessage)))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	return c.Client.Do(request)
}

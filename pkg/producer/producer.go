package producer

import (
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/backends/result"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
)

type Producer struct {
	Taskserver *machinery.Server
}

func NewProducerWithConfig(c *config.Config) (*Producer, error) {
	taskserver, err := machinery.NewServer(c)

	if err != nil {
		return nil, err
	}
	return &Producer{Taskserver: taskserver}, nil
}

func NewProducer() (*Producer, error) {
	taskserver, err := machinery.NewServer(&config.Config{
		Broker:        "redis://localhost:6379",
		ResultBackend: "redis://localhost:6379",
	})

	if err != nil {
		return nil, err
	}
	return &Producer{Taskserver: taskserver}, nil
}

func (p *Producer) Produce(contents string) (*result.AsyncResult, error) {
	task := tasks.Signature{
		Name: "handyman",
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: contents,
			},
		},
	}

	res, err := p.Taskserver.SendTask(&task)
	if err != nil {
		return nil, err
	}

	return res, nil
}

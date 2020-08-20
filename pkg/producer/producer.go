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

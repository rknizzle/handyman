package producer

import (
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
)

type Producer struct {
	Taskserver *machinery.Server
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

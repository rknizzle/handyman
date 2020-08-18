package consumer

import (
	"github.com/RichardKnop/machinery/v1/config"
	"gopkg.in/yaml.v2"
)

// ConsumerConfig contains the machinery config used for setting
// values for dealing with brokers/backends/tasks as well as
// some values specific to the handyman Consumer such as the URL
// of the application server and the number of tasks to run
// concurrently
type ConsumerConfig struct {
	MachineryConfig *config.Config
	AppURL          string `yaml:"appURL"`
	Concurrency     int    `yaml:"concurrency"`
}

func NewConsumerConfig(cfgPath string) (*ConsumerConfig, error) {
	cfg := &ConsumerConfig{}

	machineryCfg, err := config.NewFromYaml(cfgPath, false)
	if err != nil {
		return nil, err
	}
	cfg.MachineryConfig = machineryCfg

	data, err := config.ReadFromFile(cfgPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

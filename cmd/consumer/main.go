package main

import (
	"github.com/rknizzle/handyman/pkg/consumer"
)

func main() {
	c, err := consumer.NewConsumerFromCfgFile("config.yml")
	if err != nil {
		panic(err)
	}

	err = c.Start()
	if err != nil {
		panic(err)
	}
}

package main

import (
	"log"

	"github.com/southworks/gnalog/demo/auditory/config"
	"github.com/southworks/gnalog/demo/auditory/kafka"
)

func main() {
	c, err := config.Get()
	if err != nil {
		panic(err)
	}
	consumer := kafka.Consumer{}
	log.Println("Connecting to auditory...")
	consumer.Connect(c.Kafka.Brokers, "1", "logs")
	log.Println("Auditory connected")
	defer consumer.Close()
	for {
		err = consumer.Consume()
		if err != nil {
			log.Println("ERR on consume: ", err)
		}
	}
}

package config

import (
	"encoding/json"
	"os"
)

type kafkaConfig struct {
	Brokers []string
}

type conf struct {
	Kafka kafkaConfig
}

func Get() (conf, error) {
	file, _ := os.Open("./conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	c := conf{}
	err := decoder.Decode(&c)
	return c, err
}

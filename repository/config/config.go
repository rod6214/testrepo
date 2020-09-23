package config

import (
	"encoding/json"
	"os"
)

type dynamodbConfig struct {
	Region   string
	Endpoint string
}

type postgresqlConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}

type kafkaConfig struct {
	Brokers []string
}

type conf struct {
	Dynamodb   dynamodbConfig
	Postgresql postgresqlConfig
	Kafka      kafkaConfig
}

func Get() (conf, error) {
	file, _ := os.Open("./conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	c := conf{}
	err := decoder.Decode(&c)
	return c, err
}

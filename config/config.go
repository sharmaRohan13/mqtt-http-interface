package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Mqtt struct {
	Broker   string `yaml:"broker"`
	ClientId string `yaml:"client_id"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Http struct {
	Server string `yaml:"server"`
}

type Config struct {
	Mqtt Mqtt `yaml:"mqtt"`
	Http Http `yaml:"http"`
}

func ParseConfig(file string) *Config {
	// Read the yaml file
	cfg := Config{}
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Unmarshal the yaml file into the config struct
	err = yaml.Unmarshal([]byte(data), &cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return &cfg
}

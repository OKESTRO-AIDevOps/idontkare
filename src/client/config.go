package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ClientConfig struct {
	ApixVersion int    `yaml:"apix_version"`
	ConnectUrl  string `yaml:"connect_url"`
	SocketMode  bool   `yaml:"socket_mode"`
	TimeoutMS   int    `yaml:"timeout_ms"`
}

func GetClientConfig() (*ClientConfig, error) {

	var cfg ClientConfig

	file_b, err := os.ReadFile("config.yaml")

	if err != nil {

		return nil, fmt.Errorf("failed to get client config: config.yaml: %s", err.Error())

	}

	err = yaml.Unmarshal(file_b, &cfg)

	if err != nil {

		return nil, fmt.Errorf("failed to get client config: unmarshal: %s", err.Error())
	}

	return &cfg, nil
}

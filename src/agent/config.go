package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AgentConfig struct {
	ApixVersion int    `yaml:"apix_version"`
	ConnectUrl  string `yaml:"connect_url"`
	Name        string `yaml:"name"`
	UserName    string `yaml:"username"`
	KeyPath     string `yaml:"key_path"`
	UserPass    string `yaml:"userpass"`
	TimeoutMS   int    `yaml:"timeout_ms"`
}

func GetAgentConfig() (*AgentConfig, error) {

	var cfg AgentConfig

	file_b, err := os.ReadFile("config.yaml")

	if err != nil {

		return nil, fmt.Errorf("failed to get agent config: config.yaml: %s", err.Error())

	}

	err = yaml.Unmarshal(file_b, &cfg)

	if err != nil {

		return nil, fmt.Errorf("failed to get agent config: unmarshal: %s", err.Error())
	}

	return &cfg, nil
}

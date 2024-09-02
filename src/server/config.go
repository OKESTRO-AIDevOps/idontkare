package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	ApixVersion      int    `yaml:"apix_version"`
	ListenAddr       string `yaml:"listen_addr"`
	ListenPortClient string `yaml:"listen_port_client"`
	ListenPortAgent  string `yaml:"listen_port_agent"`
	ClientPath       string `yaml:"client_path"`
	AgentPath        string `yaml:"agent_path"`
	ClientTimeoutMS  int    `yaml:"client_timeout_ms"`
	AgentTimeoutMS   int    `yaml:"agent_timeout_ms"`
	DBAddr           string `yaml:"db_addr"`
	DBName           string `yaml:"db_name"`
	DBId             string `yaml:"db_id"`
	DBPw             string `yaml:"db_pw"`
	ResetDBAtStart   bool   `yaml:"reset_db_at_start"`
}

func GetServerConfig() (*ServerConfig, error) {

	var cfg ServerConfig

	file_b, err := os.ReadFile("config.yaml")

	if err != nil {

		return nil, fmt.Errorf("failed to get server config: config.yaml: %s", err.Error())

	}

	err = yaml.Unmarshal(file_b, &cfg)

	if err != nil {

		return nil, fmt.Errorf("failed to get server config: unmarshal: %s", err.Error())
	}

	return &cfg, nil
}

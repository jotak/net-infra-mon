package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	LogLevel  string `yaml:"logLevel,omitempty"`
	Server    Server `yaml:"server,omitempty"`
	PProfPort int    `yaml:"pprofPort,omitempty"`
}

type Server struct {
	Port     int    `yaml:"port,omitempty"`
	CertPath string `yaml:"certPath,omitempty"`
	KeyPath  string `yaml:"keyPath,omitempty"`
}

func Read(path string) (*Config, error) {
	// Default values
	cfg := Config{
		LogLevel: "info",
		Server: Server{
			Port: 8080,
		},
	}
	if len(path) == 0 {
		return &cfg, nil
	}
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, err
}

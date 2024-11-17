package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Main server loaded config
type Config struct {
	Server  Server  `yaml:"server"`
	Paths   Paths   `yaml:"paths"`
	Auth    Auth    `yaml:"auth"`
	Logging Logging `yaml:"logging"`
}

// Main backend web server config
type Server struct {
	Port string `yaml:"port"`
}

// File uploads etc pathing
type Paths struct {
	Public  string `yaml:"public"`
	Uploads string `yaml:"uploads"`
}

type Auth struct {
	Token string `yaml:"token"`
}

type Logging struct {
	Level string `yaml:"level"`
}

/*
Load up the config file and check for errors
*/
func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fmt.Printf("Loaded config: '%s'\n", file.Name())

	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

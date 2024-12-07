package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/txj-xyz/neoserve/file-server/internal/logger"
	"gopkg.in/yaml.v3"
)

var (
	cfg     *Config
	cfgOnce sync.Once
	log     *logger.Logger
)


/*
Load up the config file and check for errors
*/
func LoadConfig(filePath string) (*Config, error) {
	log = logger.New()

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var conf Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&conf); err != nil {
		return nil, err
	}

	SetConfig(&conf)

	// Print the loadded up config in dev mode
	if conf.Server.DevMode {
		// log.Debug("-- Config Loaded [%p] --\n%s", &conf, conf.Json())
		log.Info(
			"config loaded",
			"ptr", fmt.Sprintf("%p", &conf),
		)
		fmt.Printf("%v", conf.YAML())
	}

	return &conf, nil
}

// Set the global variable of the config
func SetConfig(config *Config) {
	cfgOnce.Do(func() { cfg = config })
}

// Fetch the latest config
func GetConfig() (*Config, error) {
	if cfg.Server.DevMode {
		// fmt.Printf("Retreiving config: [%p]\n", &cfg)
		log.Info(
			"grabbing loaded config",
			"status", fmt.Sprintf("%p", &cfg),
		)
	}

	if cfg.YAML() == "{}" {
		return nil, fmt.Errorf("failed to load config")
	} else {
		return cfg, nil
	}
}

// Pretty print JSON
func (c *Config) YAML() string {
	pretty, err := yaml.Marshal(c)
	if err != nil {
		return "{}"
	}
	return string(pretty)
}

func (p *Paths) UploadsPath() string {
	pwd, _ := os.Getwd()
	return filepath.Join(pwd, p.Uploads)
}

func (p *Paths) PublicPath() string {
	pwd, _ := os.Getwd()
	return filepath.Join(pwd, p.Public)
}

// Generate a FQDN for the instance running
func (c *Server) GenerateURL() string {
	if c.DevMode {
		return fmt.Sprintf("http://localhost:%d", c.Port)
	}

	if c.Port == 443 {
		return fmt.Sprintf("https://%s", c.Domain)
	}

	return fmt.Sprintf("http://%s:%d", c.Domain, c.Port)
}

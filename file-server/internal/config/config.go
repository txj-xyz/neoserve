package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	cfg     *Config
	cfgOnce sync.Once
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
	Port    string `yaml:"port"`
	DevMode bool   `yaml:"devMode"`
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

	var conf Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&conf); err != nil {
		return nil, err
	}

	SetConfig(&conf)

	// Print the loadded up config in dev mode
	if conf.Server.DevMode {
		fmt.Printf("Loaded config: '%+v'\n", conf)
	}

	return &conf, nil
}

// Set the global variable of the config
func SetConfig(config *Config) {
	cfgOnce.Do(func() { cfg = config })
}

func GetConfig() *Config {
	return cfg
}

// Pretty print JSON
func (c *Config) Json() string {
	fmt.Printf("%+v\n", c)
	return fmt.Sprintf("%+v\n", c)
}

// Print out a directory structure from a path
func PrintDirTree(path string, indent string) {
	// Open the directory
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for i, file := range files {
		// Print the current file or directory
		if i == len(files)-1 {
			fmt.Printf("%s└── %s\n", indent, file.Name())
		} else {
			fmt.Printf("%s├── %s\n", indent, file.Name())
		}

		// If it's a directory, recursively print its structure
		if file.IsDir() {
			PrintDirTree(filepath.Join(path, file.Name()), indent+"│   ")
		}
	}
}

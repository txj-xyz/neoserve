package config

// Root Level
type Config struct {
	Server  Server  `yaml:"server"`
	Paths   Paths   `yaml:"paths"`
	Auth    Auth    `yaml:"auth"`
	Logging Logging `yaml:"logging"`
}

type Server struct {
	Port    int    `yaml:"port"`
	MaxUpload int64 `yaml:"max_upload"`
	DevMode bool   `yaml:"dev_mode"`
	Domain  string `yaml:"domain"`
}

type Paths struct {
	Public  string `yaml:"public"`
	Uploads string `yaml:"uploads"`
}

type Auth struct {
	Token string `yaml:"token"`
}

type Logging struct {
	Discord DiscordLogModule `yaml:"discord"`
}



// Modules

type DiscordLogModule struct {
	Enabled    bool   `yaml:"enabled"`
	WebhookURL string `yaml:"webhook_url"`
}
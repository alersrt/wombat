package internal

import (
	"gopkg.in/yaml.v2"
	"log/slog"
	"mvdan.cc/sh/v3/shell"
	"os"
)

type Bot struct {
	Tag string `yaml:"tag,omitempty"`
}

type Jira struct {
	Url string `yaml:"url,omitempty"`
}

type Telegram struct {
	Token string `yaml:"token,omitempty"`
}

type PostgreSQL struct {
	Url string `yaml:"url"`
}

type Database struct {
	PostgreSQL *PostgreSQL `yaml:"postgresql,omitempty"`
}

type Config struct {
	isInitiated bool
	*Bot        `yaml:"bot,omitempty"`
	*Jira       `yaml:"jira,omitempty"`
	*Telegram   `yaml:"telegram,omitempty"`
	*Database   `yaml:"database,omitempty"`
}

func (c *Config) Init(args []string) {
	slog.Info("Wombat initialization...")

	configPath := args[0]

	defer slog.Info("Config file: " + configPath)

	file, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	replaced, err := shell.Expand(string(file), nil)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal([]byte(replaced), c)
	if err != nil {
		panic(err)
	}

	c.isInitiated = true
	return
}

func (c *Config) IsInitiated() bool {
	return c.isInitiated
}

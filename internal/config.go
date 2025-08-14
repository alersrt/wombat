package internal

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"log/slog"
	"mvdan.cc/sh/v3/shell"
	"os"
	"wombat/pkg"
)

type Bot struct {
	Tag string `yaml:"tag,omitempty"`
}

type Cipher struct {
	Key string `yaml:"key"`
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
	*Cipher     `yaml:"cipher,omitempty"`
	*Jira       `yaml:"jira,omitempty"`
	*Telegram   `yaml:"telegram,omitempty"`
	*Database   `yaml:"database,omitempty"`
}

var _ pkg.Config = (*Config)(nil)

func (c *Config) Init(args []string) error {
	slog.Info("Wombat initialization...")
	configPath := args[0]

	defer slog.Info("Config file: " + configPath)

	file, err := os.ReadFile(configPath)
	if err != nil {
		return errors.WithStack(err)
	}

	replaced, err := shell.Expand(string(file), nil)
	if err != nil {
		return errors.WithStack(err)
	}
	err = yaml.Unmarshal([]byte(replaced), c)
	if err != nil {
		return errors.WithStack(err)
	}
	c.isInitiated = true
	return nil
}

func (c *Config) IsInitiated() bool {
	return c.isInitiated
}

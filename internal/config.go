package internal

import (
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

func (c *Config) Init(args []string) (err error) {
	slog.Info("Wombat initialization...")
	configPath := args[0]

	defer slog.Info("Config file: " + configPath)
	defer pkg.CatchWithReturn(&err)

	file, err := os.ReadFile(configPath)
	pkg.Throw(err)

	replaced, err := shell.Expand(string(file), nil)
	pkg.Throw(err)

	err = yaml.Unmarshal([]byte(replaced), c)
	pkg.Throw(err)

	c.isInitiated = true
	return
}

func (c *Config) IsInitiated() bool {
	return c.isInitiated
}

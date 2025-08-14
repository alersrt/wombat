package internal

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"log/slog"
	"mvdan.cc/sh/v3/shell"
	"os"
	"sync"
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
	mtx      sync.Mutex
	Bot      `yaml:"bot,omitempty"`
	Cipher   `yaml:"cipher,omitempty"`
	Jira     `yaml:"jira,omitempty"`
	Telegram `yaml:"telegram,omitempty"`
	Database `yaml:"database,omitempty"`
}

var _ pkg.Config = (*Config)(nil)

func (c *Config) Init(args []string) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	slog.Info("Wombat initialization...")
	configPath := args[0]

	defer slog.Info("Config file: " + configPath)

	file, err := os.ReadFile(configPath)
	if err != nil {
		return errors.New(err.Error())
	}

	replaced, err := shell.Expand(string(file), nil)
	if err != nil {
		return errors.New(err.Error())
	}
	err = yaml.Unmarshal([]byte(replaced), c)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

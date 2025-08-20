package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"log/slog"
	"mvdan.cc/sh/v3/shell"
	"os"
	"sync"
)

var (
	ErrConfig     = errors.New("config")
	ErrConfigInit = errors.New("config: init")
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

func (c *Config) Init(args []string) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	slog.Info("config: init: start")
	defer slog.Info("config: init: finish")

	configPath := args[0]

	defer slog.Info("config: file: " + configPath)

	file, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrConfigInit, err)
	}

	replaced, err := shell.Expand(string(file), nil)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrConfigInit, err)
	}
	err = yaml.Unmarshal([]byte(replaced), c)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrConfigInit, err)
	}
	return nil
}

package app

import (
	"gopkg.in/yaml.v2"
	"log/slog"
	"mvdan.cc/sh/v3/shell"
	"os"
)

type Bot struct {
	Tag   string `yaml:"tag,omitempty"`
	Emoji string `yaml:"emoji,omitempty"`
}

type Jira struct {
	Url      string `yaml:"url,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

type Telegram struct {
	Token   string `yaml:"token,omitempty"`
	Webhook string `yaml:"webhook,omitempty"`
}

type Kafka struct {
	GroupId   string `yaml:"group-id,omitempty"`
	Bootstrap string `yaml:"bootstrap,omitempty"`
	Topic     string `yaml:"topic,omitempty"`
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
	*Kafka      `yaml:"kafka,omitempty"`
	*Database   `yaml:"database,omitempty"`
}

func (receiver *Config) Init(args []string) error {
	slog.Info("Wombat initialization...")

	configPath := args[0]

	defer slog.Info("Config file: " + configPath)

	file, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	replaced, err := shell.Expand(string(file), nil)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal([]byte(replaced), receiver)
	if err != nil {
		return err
	}

	if err == nil {
		receiver.isInitiated = true
	}
	return nil
}

func (receiver *Config) IsInitiated() bool {
	return receiver.isInitiated
}

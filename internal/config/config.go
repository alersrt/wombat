package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log/slog"
	"mvdan.cc/sh/v3/shell"
	"os"
)

type Bot struct {
	Tag   string `yaml:"tag,omitempty"`
	Emoji string `yaml:"emoji,omitempty"`
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
	Host     string `yaml:"host,omitempty"`
	Port     int    `yaml:"port,omitempty"`
	Database string `yaml:"database,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

func (receiver *PostgreSQL) FormatURL() string {
	return fmt.Sprintf("postgres://%[1]s:%[2]s@%[3]s:%[4]d/%[5]s",
		receiver.Username, receiver.Password, receiver.Host, receiver.Port, receiver.Database)
}

type Database struct {
	PostgreSQL *PostgreSQL `yaml:"postgresql,omitempty"`
}

type Config struct {
	isInitiated bool
	*Bot        `yaml:"bot,omitempty"`
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

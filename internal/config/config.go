package config

import (
	"gopkg.in/yaml.v2"
	"log/slog"
	"os"
)

type Bot struct {
	Tag   string `yaml:"tag"`
	Emoji string `yaml:"emoji"`
}

type Telegram struct {
	Token   string `yaml:"token"`
	Webhook string `yaml:"webhook"`
}

type Kafka struct {
	GroupId   string `yaml:"group-id"`
	Bootstrap string `yaml:"bootstrap"`
	Topic     string `yaml:"topic"`
}

type Config struct {
	isInitiated bool
	*Bot        `yaml:"bot"`
	*Telegram   `yaml:"telegram"`
	*Kafka      `yaml:"kafka"`
}

func (receiver *Config) Init(args []string) error {
	slog.Info("Wombat initialization...")

	configPath := args[0]

	defer slog.Info("Config file: " + configPath)

	file, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	replaced := os.ExpandEnv(string(file))
	err = yaml.Unmarshal([]byte(replaced), receiver)

	if err == nil {
		receiver.isInitiated = true
	}

	return err
}

func (receiver *Config) IsInitiated() bool {
	return receiver.isInitiated
}

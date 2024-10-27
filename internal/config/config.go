package config

import (
	"flag"
	"gopkg.in/yaml.v2"
	"log/slog"
	"os"
)

type Config struct {
	isInitiated bool
	Bot         struct {
		Tag   string `yaml:"tag"`
		Emoji string `yaml:"emoji"`
	} `yaml:"bot"`
	Telegram struct {
		Token   string `yaml:"token"`
		Webhook string `yaml:"webhook"`
	} `yaml:"telegram"`
	Kafka struct {
		GroupId   string `yaml:"group-id"`
		Bootstrap string `yaml:"bootstrap"`
		Topic     string `yaml:"topic"`
	} `yaml:"kafka"`
}

func (receiver *Config) Init(args []string) error {
	slog.Info("Wombat initialization...")

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	configPath := flags.String("config", "./cmd/config.yaml", "path to config")

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	defer slog.Info("Config file: " + *configPath)

	file, err := os.ReadFile(*configPath)
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

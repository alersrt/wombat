package config

import (
	"flag"
	"gopkg.in/yaml.v2"
	"os"
	. "wombat/pkg/log"
)

type Config struct {
	Bot struct {
		Tag   string `yaml:"tag"`
		Emoji string `yaml:"emoji"`
	} `yaml:"bot"`
	Telegram struct {
		Token   string `yaml:"token"`
		Webhook string `yaml:"webhook"`
	} `yaml:"telegram"`
	Kafka struct {
		ClientId  string `yaml:"client-id"`
		Bootstrap string `yaml:"bootstrap"`
		Topic     string `yaml:"topic"`
	} `yaml:"kafka"`
}

func (receiver *Config) Init(args []string) error {
	InfoLog.Print("Wombat initialization...")

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	configPath := flags.String("config", "./cmd/config.yaml", "path to config")

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	defer InfoLog.Print("Config file: " + *configPath)

	file, err := os.ReadFile(*configPath)
	if err != nil {
		return err
	}

	replaced := os.ExpandEnv(string(file))
	err = yaml.Unmarshal([]byte(replaced), receiver)

	return err
}

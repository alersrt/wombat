package config

import (
	"flag"
	"gopkg.in/yaml.v2"
	"os"
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
}

func (receiver *Config) Init(args []string) error {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	configPath := flags.String("c", "./cmd/config.yaml", "path to config")

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	config := &Config{}

	file, err := os.Open(*configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return err
	}

	return err
}

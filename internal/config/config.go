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
}

func (receiver *Config) Init(args []string) error {
	InfoLog.Println("Wombat initialization...")

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	configPath := flags.String("config", "./cmd/config.yaml", "path to config")

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	defer InfoLog.Println("Config file: " + *configPath)

	file, err := os.Open(*configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(receiver); err != nil {
		return err
	}

	return err
}

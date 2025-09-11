package internal

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"mvdan.cc/sh/v3/shell"
	"os"
	"time"
)

type Config struct {
	Imap struct {
		Address     string        `yaml:"address"`
		Username    string        `yaml:"username"`
		Password    string        `yaml:"password"`
		Mailbox     string        `yaml:"mailbox"`
		IdleTimeout time.Duration `yaml:"idleTimeout"`
	} `yaml:"imap"`
}

func NewConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: new: %v", err)
	}

	replaced, err := shell.Expand(string(file), nil)
	if err != nil {
		return nil, fmt.Errorf("config: new: %v", err)
	}

	var c Config
	err = yaml.Unmarshal([]byte(replaced), &c)
	if err != nil {
		return nil, fmt.Errorf("config: new: %v", err)
	}
	return &c, nil
}

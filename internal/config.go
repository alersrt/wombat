package internal

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"mvdan.cc/sh/v3/shell"
	"os"
	"time"
)

type Jira struct {
	Url   string `yaml:"url"`
	Token string `yaml:"token"`
}

type Telegram struct {
	Token string `yaml:"token"`
}

type Imap struct {
	Url         string        `yaml:"url"`
	Username    string        `yaml:"username"`
	Password    string        `yaml:"password"`
	Mailbox     string        `yaml:"mailbox"`
	IdleTimeout time.Duration `yaml:"idleTimeout"`
	Verbose     bool          `yaml:"verbose"`
}

type Plugin struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

type Rule struct {
	Name   string `yaml:"name"`
	Filter string `yaml:"filter"`
}

type Config struct {
	Jira     *Jira     `yaml:"jira"`
	Telegram *Telegram `yaml:"telegram"`
	Imap     *Imap     `yaml:"imap"`
	Plugins  []*Plugin `yaml:"plugins"`
	Rules    []*Rule   `yaml:"rules"`
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

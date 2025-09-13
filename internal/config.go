package internal

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"mvdan.cc/sh/v3/shell"
	"os"
)

type PluginCfg struct {
	Name string         `yaml:"name"`
	Bin  string         `yaml:"bin"`
	Conf map[string]any `yaml:"conf"`
}

type RuleCfg struct {
	Name   string `yaml:"name"`
	Filter string `yaml:"filter"`
}

type Config struct {
	Plugins []*PluginCfg `yaml:"plugins"`
	Rules   []*RuleCfg   `yaml:"rules"`
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

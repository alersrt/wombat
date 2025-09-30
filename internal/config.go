package internal

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
	"mvdan.cc/sh/v3/shell"
)

type PluginCfg struct {
	Name string `yaml:"name"`
	Bin  string `yaml:"bin"`
}

type ItemCfg struct {
	Name   string         `yaml:"name"`
	Plugin string         `yaml:"plugin"`
	Conf   map[string]any `yaml:"conf"`
}

type ApplienceCfg struct {
	Filter    string `yaml:"filter"`
	Transform string `yaml:"transform"`
}

type RuleCfg struct {
	Name     string          `yaml:"name"`
	Producer string          `yaml:"producer"`
	Consumer string          `yaml:"consumer"`
	Applies  []*ApplienceCfg `yaml:"applies"`
}

type Config struct {
	Plugins   []*PluginCfg `yaml:"plugins"`
	Producers []*ItemCfg   `yaml:"producers"`
	Consumers []*ItemCfg   `yaml:"consumers"`
	Rules     []*RuleCfg   `yaml:"rules"`
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

package internal

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
	"mvdan.cc/sh/v3/shell"
)

// ComponentType represents the type of a pipeline component
type ComponentType string

const (
	ComponentTypeSource    ComponentType = "source"
	ComponentTypeProcessor ComponentType = "processor"
	ComponentTypeSink      ComponentType = "sink"
	ComponentTypeLookup    ComponentType = "lookup"
)

// PluginCfg represents the plugin configuration for a component
type PluginCfg struct {
	Path          string            `yaml:"path,omitempty"`
	RemoteAddress string            `yaml:"remote_address,omitempty"`
	Env           map[string]string `yaml:"env,omitempty"`
}

// ExprConfig represents expression configuration
type ExprConfig struct {
	Input  string `yaml:"input,omitempty"`
	Output string `yaml:"output,omitempty"`
}

// ComponentConfig represents component-specific configuration
type ComponentConfig struct {
	Conf map[string]any `yaml:",inline"`
}

// ComponentCfg represents a component in the pipeline
type ComponentCfg struct {
	ID     string           `yaml:"id"`
	Type   ComponentType    `yaml:"type"`
	Plugin *PluginCfg       `yaml:"plugin"`
	Config *ComponentConfig `yaml:"config,omitempty"`
	Expr   *ExprConfig      `yaml:"expr,omitempty"`
}

// RoutingCfg represents a routing rule between components
type RoutingCfg struct {
	From    string `yaml:"from"`
	To      string `yaml:"to"`
	WhenCel string `yaml:"when_cel,omitempty"`
    ExprCel string `yaml:"when_cel,omitempty"`
}

// Config represents the pipeline configuration
type Config struct {
	Version    string          `yaml:"version"`
	PipelineID string          `yaml:"pipeline_id"`
	Components []*ComponentCfg `yaml:"components"`
	Routing    []*RoutingCfg   `yaml:"routing"`
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

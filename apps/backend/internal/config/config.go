package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const defaultConfigPath = "/etc/cloudflared/config.yml"

// Config represents the cloudflared configuration structure
type Config struct {
	TunnelID   string                 `yaml:"tunnel,omitempty"`
	Ingress    []IngressRule          `yaml:"ingress,omitempty"`
	Metrics    string                 `yaml:"metrics,omitempty"`
	LogLevel   string                 `yaml:"loglevel,omitempty"`
	Additional map[string]interface{} `yaml:",inline"`
}

// IngressRule represents an ingress rule in cloudflared config
type IngressRule struct {
	Hostname string `yaml:"hostname,omitempty"`
	Service  string `yaml:"service"`
	Path     string `yaml:"path,omitempty"`
}

// ReadConfig reads the cloudflared configuration from the specified path
func ReadConfig(path string) (*Config, error) {
	if path == "" {
		path = defaultConfigPath
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

// WriteConfig writes the cloudflared configuration to the specified path
func WriteConfig(path string, config *Config) error {
	if path == "" {
		path = defaultConfigPath
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write with appropriate permissions
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// ValidateConfig performs basic validation on the configuration
func ValidateConfig(config *Config) error {
	if len(config.Ingress) == 0 {
		return fmt.Errorf("ingress rules are required")
	}

	// Check that last ingress rule is a catch-all
	lastRule := config.Ingress[len(config.Ingress)-1]
	if lastRule.Service == "" {
		return fmt.Errorf("last ingress rule must have a service")
	}

	return nil
}


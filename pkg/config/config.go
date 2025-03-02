package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Server      ServerConfig      `yaml:"server"`
	Connection  ConnectionConfig  `yaml:"connection"`
	Logging     LoggingConfig     `yaml:"logging"`
	Proxy       ProxyConfig       `yaml:"proxy"`
}

// ServerConfig represents the server configuration
type ServerConfig struct {
	Host           string        `yaml:"host"`
	Port           int           `yaml:"port"`
}

// ConnectionConfig represents the connection configuration
type ConnectionConfig struct {
	BufferSize      int           `yaml:"buffer_size"`
	KeepAlive       bool          `yaml:"keep_alive"`
	KeepAlivePeriod time.Duration `yaml:"keep_alive_period"`
}

// LoggingConfig represents the logging configuration
type LoggingConfig struct {
	Level    string `yaml:"level"`
	Format   string `yaml:"format"`
	Output   string `yaml:"output"`
	FilePath string `yaml:"file_path,omitempty"` // Optional: only used when Output is "file"
}

// ProxyConfig represents the proxy configuration
type ProxyConfig struct {
	Backends []string `yaml:"backends"`
}

// Load loads the configuration from a file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
} 
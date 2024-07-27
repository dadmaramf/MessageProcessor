// Package config provides configuration loading and access functionalities.
// It reads configuration from a YAML file and makes it accessible through
// various getter methods.
package config

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"os"
)

type ConfigInterface interface {
	GetHTTPClient() HTTPClient
	GetPostgresConnect() PostgresConnect
	GetKafka() Kafka
}

type Config struct {
	Http HTTPClient      `yaml:"httpClient"`      // HTTP client configuration
	Pg   PostgresConnect `yaml:"postgresConnect"` // PostgreSQL connection
	Kf   Kafka           `yaml:"kafka"`           // Kafka configuration
}

type HTTPClient struct {
	Host string `yaml:"host"` // Host for the HTTP client
	Port string `yaml:"port"` // Port for the HTTP client
}

type PostgresConnect struct {
	Host     string `yaml:"host"`     // Host for PostgreSQL
	Port     string `yaml:"port"`     // Port for PostgreSQL
	User     string `yaml:"user"`     // User for connecting to PostgreSQL
	Password string `yaml:"password"` // Password for connecting to PostgreSQL
	DBname   string `yaml:"dbname"`   // Name of the PostgreSQL database
}

type Kafka struct {
	Brockers []string `yaml:"brockers"` // List of Kafka brokers
	Topic    string   `yaml:"topic"`    // Kafka topic
}

func (c *Config) GetHTTPClient() HTTPClient {
	return c.Http
}

func (c *Config) GetPostgresConnect() PostgresConnect {
	return c.Pg
}

func (c *Config) GetKafka() Kafka {
	return c.Kf
}

// NewConfig loads the configuration from a YAML file, the path to which is specified by the PATHCONF environment variable.
func NewConfig() (*Config, error) {
	configPath := os.Getenv("PATHCONF")
	if configPath == "" {
		return nil, fmt.Errorf("PATHCONF is not set")
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("not read yaml file %s ", err)
	}

	var config Config

	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("not convert yaml %s", err)
	}

	return &config, nil
}

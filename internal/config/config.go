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
	Http HTTPClient      `yaml:"httpClient"`
	Pg   PostgresConnect `yaml:"postgresConnect"`
	Kf   Kafka           `yaml:"kafka"`
}

type HTTPClient struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type PostgresConnect struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBname   string `yaml:"dbname"`
}

type Kafka struct {
	Brockers []string `yaml:"brockers"`
	Topic    string   `yaml:"topic"`
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

func NewConfig() (*Config, error) {

	file, err := os.ReadFile(os.Getenv("PATHCONF"))

	if err != nil {

		return nil, fmt.Errorf("not read yaml file %s ", err)
	}

	var config Config

	if err := yaml.Unmarshal(file, &config); err != nil {

		return nil, fmt.Errorf("not convert yaml %s", err)
	}

	return &config, nil
}

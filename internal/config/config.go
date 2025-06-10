package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Port     int      `yaml:"port"`
	Strategy string   `yaml:"strategy"`
	Backends []string `yaml:"backends"`
}

func GetConfig() (*Config, error) {
	f, err := os.Open("config.yaml")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) GetStrategy() string {
	if c.Strategy == "" {
		return "round-robin"
	}
	return c.Strategy
}

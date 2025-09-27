package config

import (
	"fmt"
	"gitlab.com/kforge/kforge-sdk/sdkconfig"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Postgres sdkconfig.Postgres `yaml:"postgres"`
	Http     HttpConfig         `yaml:"http"`
}

type HttpConfig struct {
	Addr string `yaml:"addr"`
}

func ParseConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %v", err)
	}
	defer f.Close()

	var cfg Config
	err = yaml.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	return &cfg, nil
}

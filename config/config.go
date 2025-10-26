package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Postgres Postgres    `yaml:"postgres"`
	Http     HttpConfig  `yaml:"http"`
	Minio    MinioConfig `yaml:"minio"`
	YooKassa YooKassa   `yaml:"yookassa"`
}

type HttpConfig struct {
	Addr string `yaml:"addr"`
}

type MinioConfig struct {
	Endpoint   string `yaml:"endpoint"`
	AccessKey  string `yaml:"access_key"`
	SecretKey  string `yaml:"secret_key"`
	UseSSL     bool   `yaml:"use_ssl"`
	BucketName string `yaml:"bucket_name"`
	MaxRetries int    `yaml:"max_retries"`
}

const dsnTemplate = "host=%s port=%s user=%s password=%s dbname=%s application_name=%s sslmode=disable"

type Postgres struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	AppName  string `yaml:"app_name"`
	MaxConns int    `yaml:"max_conns"`
}

func (p *Postgres) MustValidate() {
	err := p.Validate()
	if err != nil {
		panic(err.Error())
	}
}

func (p *Postgres) Validate() error {
	if p.Host == "" {
		return fmt.Errorf("host is required")
	}
	if p.Port == "" {
		return fmt.Errorf("port is required")
	}
	if p.Username == "" {
		return fmt.Errorf("username is required")
	}
	if p.Password == "" {
		return fmt.Errorf("password is required")
	}
	if p.Database == "" {
		return fmt.Errorf("database is required")
	}
	if p.AppName == "" {
		return fmt.Errorf("app_name is required")
	}
	if p.MaxConns == 0 {
		return fmt.Errorf("max_conns is required")
	}

	return nil
}

func (p *Postgres) ToDSN() string {
	return fmt.Sprintf(dsnTemplate,
		p.Host,
		p.Port,
		p.Username,
		p.Password,
		p.Database,
		p.AppName,
	)
}

type YooKassa struct {
	Host      string `yaml:"host"`
	SecretKey string `yaml:"secret_key"`
	AccountID string `yaml:"account_id"`
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

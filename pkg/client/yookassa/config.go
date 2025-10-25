package yookassa

import "net/url"

type Config struct {
	Host      *url.URL `yaml:"host"`
	SecretKey string   `yaml:"secret_key"`
	AccountID string   `yaml:"account_id"`
}

package yookassa

import (
	"github.com/go-resty/resty/v2"
	"net/url"
)

type Client struct {
	client *resty.Client

	host      url.URL
	secretKey string
	accountID string
}

func New(
	client *resty.Client,
	cfg Config,
) *Client {
	return &Client{
		client:    client,
		host:      cfg.Host,
		secretKey: cfg.SecretKey,
		accountID: cfg.AccountID,
	}
}

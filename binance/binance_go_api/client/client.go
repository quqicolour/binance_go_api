package client

import (
	"binance_connector"
	"time"
)

type Client struct {
	Conn      *binance_connector.Client
	APIKey    string
	SecretKey string
	Timeout   time.Duration
	BaseAPI   string
	BaseWS    string
	ProxyURL  string
}

func NewClient(binanceClient *binance_connector.Client, apiKey, secretKey, baseAPI, baseWS, proxyURL string) *Client {
	return &Client{
		Conn:      binanceClient,
		APIKey:    apiKey,
		SecretKey: secretKey,
		Timeout:   time.Second * 15,
		BaseAPI:   baseAPI,
		BaseWS:    baseWS,
		ProxyURL:  proxyURL,
	}
}

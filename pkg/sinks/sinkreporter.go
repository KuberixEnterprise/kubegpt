package sinks

import (
	"net/http"
	"time"
)

type Client struct {
	hclient *http.Client
}

func NewClient() *Client {
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	return &Client{
		hclient: client,
	}
}

package novakeyclient

import (
	"net/http"	
	"strings"	 
)

func NewClient() *Client {
	cfg := GetConfig()
	return &Client{
		BaseURL:    strings.TrimRight(cfg.Endpoint, "/"),
		HTTPClient: http.DefaultClient,
	}
}
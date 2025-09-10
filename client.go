package novakeyauth

import (
	"net/http"	
	"strings"	
  "novakeyauth/internal/client"
)

func NewClient(baseURL string) *client.Client {
	return &client.Client{
		BaseURL:    strings.TrimRight(baseURL, "/"),
		HTTPClient: http.DefaultClient,
	}
}
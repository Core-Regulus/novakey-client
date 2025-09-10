package client

import (
	"context"
	cryptoRand "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"novakeyauth/internal/keys"
	"novakeyauth/internal/signedRequest"
	"novakeyauth/internal/users"
	"novakeyauth/internal/workspaces"
	"strings"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func GenerateNonce(n int) ([]byte, error) {
	if n <= 0 {
		n = 32
	}
	b := make([]byte, n)
	if _, err := cryptoRand.Read(b); err != nil {
		return nil, fmt.Errorf("rand: %w", err)
	}
	return b, nil
}


func sign(req *signedrequest.SignedRequest, privateKey string) error {    
  messageBytes, err := GenerateNonce(32)
	if err != nil {
		return fmt.Errorf("generate nonce: %w", err)
	}
	
  signer, pub, err := keys.ParseOpenSSHED25519Signer(privateKey)
  if err != nil {
	  return fmt.Errorf("parse private key: %w", err)
  }
	  
  sig, err := signer.Sign(cryptoRand.Reader, messageBytes)
	if (err != nil) {
	  return fmt.Errorf("signer error: %w", err)
  }
  req.PublicKey, err = keys.EncodeSSHPublicKey(pub, req.Username)
  if (err != nil) {
	  return fmt.Errorf("encode public key: %w", err)
  }

	req.Signature = base64.StdEncoding.EncodeToString(sig.Blob)
	req.Message = base64.StdEncoding.EncodeToString(messageBytes)
	if req.Timestamp == 0 {
		req.Timestamp = time.Now().Unix()
	}
  return nil
}

func send[Request any, Response any](c *Client, ctx context.Context, req Request, relativeUrl string) (*Response, *http.Response, error) {
  body, err := json.Marshal(req)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal request: %w", err)
	}

	url := c.BaseURL + relativeUrl
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(body)))
	if err != nil {
		return nil, nil, fmt.Errorf("build request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	httpClient := c.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, nil, fmt.Errorf("do request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, fmt.Errorf("read response: %w", err)
	}

	var parsed Response
	_ = json.Unmarshal(respBytes, &parsed)	
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &parsed, resp, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, strings.TrimSpace(string(respBytes)))
	}
  return &parsed, resp, nil
}

func (c *Client) NewUser (
	ctx context.Context,
	privateKey string,
	req users.SetUserRequest,
) (*users.SetUserResponse, *http.Response, error) {

	if c == nil {
		return nil, nil, errors.New("nil client")
	}

  sign(&req.SignedRequest, privateKey)  
	return send[users.SetUserRequest, users.SetUserResponse](c, ctx, req, "/users/set");
}

func (c *Client) DeleteUser (
	ctx context.Context,
	privateKey string,
	req users.DeleteUserRequest,
) (*users.DeleteUserResponse, *http.Response, error) {

	if c == nil {
		return nil, nil, errors.New("nil client")
	}

  sign(&req.SignedRequest, privateKey)
  return send[users.DeleteUserRequest, users.DeleteUserResponse](c, ctx, req, "/users/delete");
}

func (c *Client) NewWorkspace (
	ctx context.Context,
	privateKey string,
	req workspaces.SetWorkspaceRequest,
) (*workspaces.SetWorkspaceResponse, *http.Response, error) {

	if c == nil {
		return nil, nil, errors.New("nil client")
	}

  sign(&req.SignedRequest, privateKey)  
	return send[workspaces.SetWorkspaceRequest, workspaces.SetWorkspaceResponse](c, ctx, req, "/workspaces/set");
}

func (c *Client) DeleteWorkspace (
	ctx context.Context,
	privateKey string,
	req workspaces.DeleteWorkspaceRequest,
) (*workspaces.DeleteWorkspaceResponse, *http.Response, error) {

	if c == nil {
		return nil, nil, errors.New("nil client")
	}

  sign(&req.SignedRequest, privateKey)
  return send[workspaces.DeleteWorkspaceRequest, workspaces.DeleteWorkspaceResponse](c, ctx, req, "/workspaces/delete");
}
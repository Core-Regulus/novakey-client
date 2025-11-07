package novakeyclient

import (
	"context"
	cryptoRand "crypto/rand"
	"encoding/base64"
	"encoding/json"	
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"github.com/core-regulus/novakey-types-go"	
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


func sign(req *novakeytypes.AuthEntity, privateKey string) error {    
  messageBytes, err := GenerateNonce(32)
	if err != nil {
		return fmt.Errorf("generate nonce: %w", err)
	}
	
  signer, pub, err := ParseOpenSSHED25519Signer(privateKey)
  if err != nil {
	  return fmt.Errorf("parse private key: %w", err)
  }
	  
  sig, err := signer.Sign(cryptoRand.Reader, messageBytes)
	if (err != nil) {
	  return fmt.Errorf("signer error: %w", err)
  }
  req.PublicKey, err = EncodeSSHPublicKey(pub, req.Username)
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
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return &parsed, resp, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, strings.TrimSpace(string(respBytes)))
	}
  return &parsed, resp, nil
}

func (c *Client) SetUser (
	ctx context.Context,
	privateKey string,
	signerKey string,
	req novakeytypes.SetUserRequest,
) (novakeytypes.SetUserResponse) {

	if c == nil {
		return novakeytypes.SetUserResponse{
			Error: novakeytypes.Error{
				Status:  0,
				Error:   "nil client",				
			},
		}		
	}

  sign(&req.AuthEntity, privateKey)
	sign(&req.Signer, signerKey)
	uResp, resp, err := send[novakeytypes.SetUserRequest, novakeytypes.SetUserResponse](c, ctx, req, "/users/set");
	if (err != nil) {
		return novakeytypes.SetUserResponse{
			Error: novakeytypes.Error{
				Status:  resp.StatusCode,
				Error:   err.Error(),				
			},
		}		
	}
	return *uResp
}

func (c *Client) DeleteUser (
	ctx context.Context,
	privateKey string,
	req novakeytypes.DeleteUserRequest,
) (novakeytypes.DeleteUserResponse) {

	if c == nil {
		return novakeytypes.DeleteUserResponse{
			Error: novakeytypes.Error{
				Status:  0,
				Error:   "nil client",				
			},
		}		
	}

  sign(&req.AuthEntity, privateKey)
  uResp, resp, err := send[novakeytypes.DeleteUserRequest, novakeytypes.DeleteUserResponse](c, ctx, req, "/users/delete");
	if (err != nil) {
		return novakeytypes.DeleteUserResponse{
			Error: novakeytypes.Error{
				Status:  resp.StatusCode,
				Error:   err.Error(),				
			},
		}		
	}
	return *uResp
}

func (c *Client) SetWorkspace (
	ctx context.Context,
	privateKey string,
	req novakeytypes.SetWorkspaceRequest,
) (novakeytypes.SetWorkspaceResponse) {

	if c == nil {
		return novakeytypes.SetWorkspaceResponse{
			Error: novakeytypes.Error{
				Status:  0,
				Error:   "nil client",				
			},
		}		
	}

  sign(&req.Signer, privateKey)  
	uResp, resp, err := send[novakeytypes.SetWorkspaceRequest, novakeytypes.SetWorkspaceResponse](c, ctx, req, "/workspaces/set");
	if (err != nil) {
		return novakeytypes.SetWorkspaceResponse{
			Error: novakeytypes.Error{
				Status:  resp.StatusCode,
				Error:   err.Error(),				
			},
		}		
	}
	return *uResp
}

func (c *Client) DeleteWorkspace (
	ctx context.Context,	
	privateKey string,
	req novakeytypes.DeleteWorkspaceRequest,
) (novakeytypes.DeleteWorkspaceResponse) {

	if c == nil {
		return novakeytypes.DeleteWorkspaceResponse{
			Error: novakeytypes.Error{
				Status:  0,
				Error:   "nil client",				
			},
		}		
	}

  sign(&req.Signer, privateKey)
  uResp, resp, err := send[novakeytypes.DeleteWorkspaceRequest, novakeytypes.DeleteWorkspaceResponse](c, ctx, req, "/workspaces/delete");
	if (err != nil) {
		return novakeytypes.DeleteWorkspaceResponse{
			Error: novakeytypes.Error{
				Status:  resp.StatusCode,
				Error:   err.Error(),				
			},
		}		
	}
	return *uResp
}

func (c *Client) SetProject (
	ctx context.Context,
	privateKey string,	
	req novakeytypes.SetProjectRequest,
) (novakeytypes.SetProjectResponse) {

	if c == nil {
		return novakeytypes.SetProjectResponse{
			Error: novakeytypes.Error{
				Status:  0,
				Error:   "nil client",				
			},
		}		
	}

  sign(&req.Signer, privateKey)  
	uResp, resp, err := send[novakeytypes.SetProjectRequest, novakeytypes.SetProjectResponse](c, ctx, req, "/projects/set");
	if (err != nil) {
		return novakeytypes.SetProjectResponse{
			Error: novakeytypes.Error{
				Status:  resp.StatusCode,
				Error:   err.Error(),				
			},
		}		
	}
	return *uResp
}

func (c *Client) DeleteProject (
	ctx context.Context,	
	privateKey string,
	req novakeytypes.DeleteProjectRequest,
) (novakeytypes.DeleteProjectResponse) {

	if c == nil {
		return novakeytypes.DeleteProjectResponse{
			Error: novakeytypes.Error{
				Status:  0,
				Error:   "nil client",				
			},
		}		
	}

	sign(&req.Signer, privateKey)
	uResp, resp, err := send[novakeytypes.DeleteProjectRequest, novakeytypes.DeleteProjectResponse](c, ctx, req, "/projects/delete");
	if (err != nil) {
		return novakeytypes.DeleteProjectResponse{
			Error: novakeytypes.Error{
				Status:  resp.StatusCode,
				Error:   err.Error(),				
			},
		}		
	}
	return *uResp
}

func (c *Client) GetWorkspace (
	ctx context.Context,	
	privateKey string,	
	req novakeytypes.GetWorkspaceRequest,
) (novakeytypes.GetWorkspaceResponse) {

	if c == nil {
		return novakeytypes.GetWorkspaceResponse{
			Error: novakeytypes.Error{
				Status:  0,
				Error:   "nil client",				
			},
		}		
	}

  sign(&req.Signer, privateKey)  
	uResp, resp, err := send[novakeytypes.GetWorkspaceRequest, novakeytypes.GetWorkspaceResponse](c, ctx, req, "/workspaces/get");
	if (err != nil) {
		return novakeytypes.GetWorkspaceResponse{
			Error: novakeytypes.Error{
				Status:  resp.StatusCode,
				Error:   err.Error(),				
			},
		}		
	}
	return *uResp
}

func (c *Client) GetProject (
	ctx context.Context,
	privateKey string,	
	req novakeytypes.GetProjectRequest,
) (novakeytypes.GetProjectResponse) {

	if c == nil {
		return novakeytypes.GetProjectResponse{
			Error: novakeytypes.Error{
				Status:  0,
				Error:   "nil client",				
			},
		}		
	}

  sign(&req.Signer, privateKey)  
	uResp, resp, err := send[novakeytypes.GetProjectRequest, novakeytypes.GetProjectResponse](c, ctx, req, "/projects/get");
	if (err != nil) {
		return novakeytypes.GetProjectResponse{
			Error: novakeytypes.Error{
				Status:  resp.StatusCode,
				Error:   err.Error(),				
			},
		}		
	}
	return *uResp
}


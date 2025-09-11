package workspaces

import (
	"context"
	"fmt"
	"novakeyauth/internal/signedRequest"
	"novakeyauth/internal/workspaces"
	"novakeyauth/internal/client"
	"github.com/google/uuid"
)

func CreateWorkspace(client *client.Client, priv string) (*workspaces.SetWorkspaceResponse, error) {		
	req := workspaces.SetWorkspaceRequest{ 
		Email: "testuser@test.com",
		Name: "testWorkspace",
		SignedRequest: signedrequest.SignedRequest{},
	}

	resp, _, err := client.NewWorkspace(context.Background(), priv, req)
	if err != nil {
		return nil, err
	}
	
	if resp.Status != 200 {
		return nil, fmt.Errorf("expected status ok, got %d %s", resp.Status, resp.ErrorDescription)
	}

	if resp.Password == "" {
		return nil, fmt.Errorf("expected password to be set")
	}

	return resp, nil
}

func DeleteWorkspace(client *client.Client, priv string) (uuid.UUID, error) {	
	req := workspaces.DeleteWorkspaceRequest{
		SignedRequest: signedrequest.SignedRequest{},
	}
	resp, _, err := client.DeleteWorkspace(context.Background(), priv, req)
	if err != nil {
		return uuid.Nil, err
	}
	
	if resp.Status != 200 {
		return uuid.Nil, fmt.Errorf("expected status ok, got %d %s", resp.Status, resp.ErrorDescription)
	}

	return resp.Id, nil
}

func DeleteWorkspaceByPassword(client *client.Client, id uuid.UUID, password string) (uuid.UUID, error) {
	req := workspaces.DeleteWorkspaceRequest{
		Password: password,
		Id: id,
		SignedRequest: signedrequest.SignedRequest{},
	}
	resp, _, err := client.DeleteWorkspace(context.Background(), "", req)
	if err != nil {
		return uuid.Nil, err
	}
	
	if resp.Status != 200 {
		return uuid.Nil, fmt.Errorf("expected status ok, got %d %s", resp.Status, resp.ErrorDescription)
	}

	return resp.Id, nil
}
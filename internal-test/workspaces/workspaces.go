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
		Name: "testWorkspace",
		User: signedrequest.SignedRequest{},
	}

	resp, _, err := client.NewWorkspace(context.Background(), priv, req)
	if err != nil {
		return nil, err
	}
	
	if resp.Status != 200 {		
		return nil, fmt.Errorf("%s", signedrequest.FormatErrorResponse(resp.ErrorResponse))
	}

	return resp, nil
}

func DeleteWorkspace(client *client.Client, id uuid.UUID, priv string) (uuid.UUID, error) {	
	req := workspaces.DeleteWorkspaceRequest{
		Id: id,
		User: signedrequest.SignedRequest{},
	}
	resp, _, err := client.DeleteWorkspace(context.Background(), priv, req)
	if err != nil {
		return uuid.Nil, err
	}
	
	if resp.Status != 200 {
		return uuid.Nil, fmt.Errorf("%s", signedrequest.FormatErrorResponse(resp.ErrorResponse))
	}

	return resp.Id, nil
}

func DeleteWorkspaceByPassword(client *client.Client, id uuid.UUID, userId uuid.UUID, password string) (uuid.UUID, error) {
	req := workspaces.DeleteWorkspaceRequest{		
		Id: id,
		User: signedrequest.SignedRequest{
			Id: userId,
			Password: password,
		},
	}
	resp, _, err := client.DeleteWorkspace(context.Background(), "", req)
	if err != nil {
		return uuid.Nil, err
	}
	
	if resp.Status != 200 {
		return uuid.Nil, fmt.Errorf("%s", signedrequest.FormatErrorResponse(resp.ErrorResponse))
	}

	return resp.Id, nil
}
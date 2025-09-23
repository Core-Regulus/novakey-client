package test

import (
	"context"
	"fmt"
	"novakeyclient"
	"github.com/core-regulus/novakey-types-go"
	"github.com/google/uuid"
)

func CreateProject(
	client *novakeyclient.Client, 
	workspaceId uuid.UUID, 
	priv string,
	keys []novakeytypes.Key,
) (*novakeytypes.SetProjectResponse, error) {		
	req := novakeytypes.SetProjectRequest{ 
		Name: "testProject",
		Description: "Test Project Description",		
		Keys: keys,
		WorkspaceId: workspaceId,
		Signer: novakeytypes.AuthEntity{},
	}

	resp, _, err := client.NewProject(context.Background(), priv, req)
	if err != nil {
		return nil, err
	}
	
	if resp.Status != 200 {		
		return nil, fmt.Errorf("%s", novakeytypes.FormatErrorResponse(resp.ErrorResponse))
	}

	return resp, nil
}

func DeleteProject(client *novakeyclient.Client, id uuid.UUID, priv string) (uuid.UUID, error) {	
	req := novakeytypes.DeleteProjectRequest{
		Id: id,
		Signer: novakeytypes.AuthEntity{},
	}
	resp, _, err := client.DeleteProject(context.Background(), priv, req)
	if err != nil {
		return uuid.Nil, err
	}
	
	if resp.Status != 200 {
		return uuid.Nil, fmt.Errorf("%s", novakeytypes.FormatErrorResponse(resp.ErrorResponse))
	}

	return resp.Id, nil
}
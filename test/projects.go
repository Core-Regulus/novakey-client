package test

import (
	"context"
	"fmt"
	"github.com/core-regulus/novakey-client"
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

	resp := client.SetProject(context.Background(), priv, req)	

	if resp.Status != 200 {		
		return nil, fmt.Errorf("%s", novakeytypes.FormatErrorResponse(resp.Error))
	}

	return &resp, nil
}

func DeleteProject(client *novakeyclient.Client, id uuid.UUID, priv string) (uuid.UUID, error) {	
	req := novakeytypes.DeleteProjectRequest{
		Id: id,
		Signer: novakeytypes.AuthEntity{},
	}
	resp := client.DeleteProject(context.Background(), priv, req)
	
	if resp.Status != 200 {
		return uuid.Nil, fmt.Errorf("%s", novakeytypes.FormatErrorResponse(resp.Error))
	}

	return resp.Id, nil
}

func GetProject(
	client *novakeyclient.Client, 
	Id uuid.UUID, 
	priv string) (*novakeytypes.GetProjectResponse, error) {		
	req := novakeytypes.GetProjectRequest{ 
		Id: Id,
		Signer: novakeytypes.AuthEntity{},
	}

	resp := client.GetProject(context.Background(), priv, req)	
	if resp.Status != 200 {		
		return nil, fmt.Errorf("%s", novakeytypes.FormatErrorResponse(resp.Error))
	}

	return &resp, nil
}
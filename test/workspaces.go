package test

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"novakeyclient"
	"github.com/core-regulus/novakey-types-go"

)

func CreateWorkspace(client *novakeyclient.Client, priv string) (*novakeytypes.SetWorkspaceResponse, error) {		
	req := novakeytypes.SetWorkspaceRequest{ 
		Name: "testWorkspace",
		Signer: novakeytypes.AuthEntity{},
	}

	resp := client.SetWorkspace(context.Background(), priv, req)
	
	if resp.Status != 200 {		
		return nil, fmt.Errorf("%s", novakeytypes.FormatErrorResponse(resp.Error))
	}

	return &resp, nil
}

func DeleteWorkspace(client *novakeyclient.Client, id uuid.UUID, priv string) (uuid.UUID, error) {	
	req := novakeytypes.DeleteWorkspaceRequest{
		Id: id,
		Signer: novakeytypes.AuthEntity{},
	}
	resp := client.DeleteWorkspace(context.Background(), priv, req)
	if resp.Status != 200 {
		return uuid.Nil, fmt.Errorf("%s", novakeytypes.FormatErrorResponse(resp.Error))
	}

	return resp.Id, nil
}

func DeleteWorkspaceByPassword(client *novakeyclient.Client, id uuid.UUID, userId uuid.UUID, password string) (uuid.UUID, error) {
	req := novakeytypes.DeleteWorkspaceRequest{		
		Id: id,
		Signer: novakeytypes.AuthEntity{
			Id: userId,
			Password: password,
		},
	}
	resp := client.DeleteWorkspace(context.Background(), "", req)
	if resp.Status != 200 {
		return uuid.Nil, fmt.Errorf("%s", novakeytypes.FormatErrorResponse(resp.Error))
	}

	return resp.Id, nil
}
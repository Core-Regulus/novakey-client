package test

import (
	"context"
	"fmt"	
	"testing"
	"novakeyclient"
	"github.com/core-regulus/novakey-types-go"
	"github.com/google/uuid"
)

func CreateUser(
	t *testing.T, 
	client *novakeyclient.Client, 
	signer string,
	workspaces []novakeytypes.Workspace,
	projects []novakeytypes.Project,
) (*novakeytypes.SetUserResponse, string, error) {		
	priv, err := novakeyclient.GenerateKey()
	if (err != nil) {
		t.Fatal(err);
	}
	
	req := novakeytypes.SetUserRequest{ 
		Email: "testuser@test.com",
		Workspaces: workspaces,
		Projects: projects,
		AuthEntity: novakeytypes.AuthEntity {
	  	Username: "testuser",
		},
	}

	resp, _, err := client.NewUser(context.Background(), priv, signer, req)
	if err != nil {
		t.Fatal(err)		
	}
	
	if resp.Status != 200 {
		t.Fatal(novakeytypes.FormatErrorResponse(resp.ErrorResponse));
	}

	if resp.Password == "" {
		t.Fatal(fmt.Errorf("expected password to be set"))
	}

	return resp, priv, nil
}

func DeleteUser(client *novakeyclient.Client, priv string) (uuid.UUID, error) {	
	req := novakeytypes.DeleteUserRequest{
		AuthEntity: novakeytypes.AuthEntity{
      Username: "testuser",
    },
	}
	resp, _, err := client.DeleteUser(context.Background(), priv, req)
	if err != nil {
		return uuid.Nil, err
	}
	
	if resp.Status != 200 {
		return uuid.Nil, fmt.Errorf("%s", novakeytypes.FormatErrorResponse(resp.ErrorResponse))
	}

	return resp.Id, nil
}

func DeleteUserByPassword(client *novakeyclient.Client, id uuid.UUID, password string) (uuid.UUID, error) {
	req := novakeytypes.DeleteUserRequest{		
		AuthEntity: novakeytypes.AuthEntity{
			Password: password,
			Id: id,
		},
	}
	resp, _, err := client.DeleteUser(context.Background(), "", req)
	if err != nil {
		return uuid.Nil, err
	}
	
	if resp.Status != 200 {
		return uuid.Nil, fmt.Errorf("%s", novakeytypes.FormatErrorResponse(resp.ErrorResponse))
	}

	return resp.Id, nil
}
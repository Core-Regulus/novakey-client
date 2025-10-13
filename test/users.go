package test

import (
	"context"
	"fmt"	
	"os"
	"testing"
	"github.com/core-regulus/novakey-client"
	"github.com/core-regulus/novakey-types-go"
	"github.com/google/uuid"
)

func CreateUser(
	t *testing.T, 
	client *novakeyclient.Client, 
	signer string,
	workspaces []novakeytypes.Workspace,
	projects []novakeytypes.Project,
	email string,
) (*novakeytypes.SetUserResponse, string, error) {		
	priv, err := novakeyclient.GenerateKey()
	if (err != nil) {
		t.Fatal(err);
	}
	
	if email == "" {
		email = "testuser@test.com"
	}

	req := novakeytypes.SetUserRequest{ 
		Email: email,
		Workspaces: workspaces,
		Projects: projects,
		AuthEntity: novakeytypes.AuthEntity {
	  	Username: "testuser",
		},
	}

	resp := client.SetUser(context.Background(), priv, signer, req)
	if resp.Status != 200 {
		t.Fatal(novakeytypes.FormatErrorResponse(resp.Error));
	}

	if resp.Password == "" {
		t.Fatal(fmt.Errorf("expected password to be set"))
	}

	return &resp, priv, nil
}

func DeleteUser(client *novakeyclient.Client, priv string) (uuid.UUID, error) {	
	req := novakeytypes.DeleteUserRequest{
		AuthEntity: novakeytypes.AuthEntity{
      Username: "testuser",
    },
	}
	resp := client.DeleteUser(context.Background(), priv, req)
	if resp.Status != 200 {
		return uuid.Nil, fmt.Errorf("%s", novakeytypes.FormatErrorResponse(resp.Error))
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
	resp := client.DeleteUser(context.Background(), "", req)
	if resp.Status != 200 {
		return uuid.Nil, fmt.Errorf("%s", novakeytypes.FormatErrorResponse(resp.Error))
	}

	return resp.Id, nil
}

func GeneratekeyToFile(filename string) error {
	_, err := os.Stat(filename)
	if err == nil {
		return nil;
	}

	priv, err := novakeyclient.GenerateKey()
	if (err != nil) {		
		return err
	}
	err = os.WriteFile(filename, []byte(priv), 0600)
	if err != nil {
		return err
	}
	return nil	
}
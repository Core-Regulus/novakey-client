package users

import (
	"context"
	"fmt"
	"novakeyauth/internal/signedRequest"
	"novakeyauth/internal/users"
	"novakeyauth/internal/client"
	"github.com/google/uuid"
)

func CreateUser(client *client.Client, priv string) (*users.SetUserResponse, error) {		
	req := users.SetUserRequest{ 
		Email: "testuser@test.com",
		SignedRequest: signedrequest.SignedRequest{
	  	Username: "testuser",
		},
	}

	resp, _, err := client.NewUser(context.Background(), priv, req)
	if err != nil {
		return nil, err
	}
	
	if resp.Status != 200 {
		return nil, fmt.Errorf("expected status ok, got %d", resp.Status)
	}

	if resp.Password == "" {
		return nil, fmt.Errorf("expected password to be set")
	}

	return resp, nil
}

func DeleteUser(client *client.Client, priv string) (uuid.UUID, error) {	
	req := users.DeleteUserRequest{
		SignedRequest: signedrequest.SignedRequest{
      Username: "testuser",
    },
	}
	resp, _, err := client.DeleteUser(context.Background(), priv, req)
	if err != nil {
		return uuid.Nil, err
	}
	
	if resp.Status != 200 {
		return uuid.Nil, fmt.Errorf("expected status ok, got %d", resp.Status)
	}

	return resp.Id, nil
}

func DeleteUserByPassword(client *client.Client, id uuid.UUID, password string) (uuid.UUID, error) {
	req := users.DeleteUserRequest{
		Password: password,
		Id: id,
		SignedRequest: signedrequest.SignedRequest{},
	}
	resp, _, err := client.DeleteUser(context.Background(), "", req)
	if err != nil {
		return uuid.Nil, err
	}
	
	if resp.Status != 200 {
		return uuid.Nil, fmt.Errorf("expected status ok, got %d", resp.Status)
	}

	return resp.Id, nil
}
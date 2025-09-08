package novakeyauth

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"novakeyauth/internal/config"
	"novakeyauth/internal/keys"
	"testing"
	"github.com/google/uuid"
)

func generateKey(t *testing.T) string{
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	pemPriv, err := keys.PrivateKeyToOpenSSHPEM(priv)
	if err != nil {
		t.Fatal(err)
	}	
	return pemPriv
}

func createUser(priv string) (uuid.UUID, error) {	
	conf := config.Get()
	client := NewClient(conf.Endpoint)

	req := SetUserRequest{ 
		Email: "testuser@test.com",
		SignedRequest: SignedRequest{
      Username: "testuser",
    },
	}

	resp, _, err := client.NewUser(context.Background(), priv, req)
	if err != nil {
		return uuid.Nil, err
	}
	
	if resp.Status != 200 {
		return uuid.Nil, fmt.Errorf("expected status ok, got %d", resp.Status)
	}

	return resp.Id, nil
}

func deleteUser(priv string) (uuid.UUID, error) {
	conf := config.Get()
	client := NewClient(conf.Endpoint)

	req := DeleteUserRequest{
		SignedRequest: SignedRequest{
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


func TestNewUser_Success(t *testing.T) {	
	priv := generateKey(t)
	id, err := createUser(priv)
	if err != nil {
		t.Fatalf("createUser failed: %v", err)
	}
	t.Logf("created user id=%s", id)

	rId, err := deleteUser(priv)

	if rId != id {
		t.Fatalf("deleteUser failed: %v", err)
	}

	if err != nil {
		t.Fatalf("deleteUser failed: %v", err)
	}
	t.Logf("deleted user id=%s", id)
}
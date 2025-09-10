package novakeyauth

import (
	"novakeyauth/internal-test/keys"
	"novakeyauth/internal-test/users"	
	"novakeyauth/internal-test/workspaces"	
	"novakeyauth/internal/config"
	"testing"
)

//Workspace tests

func TestSetWorkspace_Success(t *testing.T) {	
	cfg := config.Get()
	client := NewClient(cfg.Endpoint)
	priv := keys.GenerateKey(t)
	resp, err := workspaces.CreateWorkspace(client, priv)
	if err != nil {
		t.Fatalf("createWorkspace failed: %v", err)
	}
	t.Logf("created workspace id=%s", resp.Id)

	rId, err := workspaces.DeleteWorkspace(client, priv)

	if rId != resp.Id {
		t.Fatalf("deleteWorkspace failed: %v", err)
	}

	if err != nil {
		t.Fatalf("deleteWorkspace failed: %v", err)
	}
	t.Logf("deleted user id=%s", resp.Id)
}


func TestDeleteWorkspaceByPassword(t *testing.T) {	
	cfg := config.Get()
	client := NewClient(cfg.Endpoint)
	priv := keys.GenerateKey(t)
	resp, err := workspaces.CreateWorkspace(client, priv)
	if err != nil {
		t.Fatalf("createWorkspace failed: %v", err)
	}
	t.Logf("created workspace id=%s", resp.Id)

	rId, err := workspaces.DeleteWorkspaceByPassword(client, resp.Id, resp.Password)

	if rId != resp.Id {
		t.Fatalf("deleteWorkspace failed: %v", err)
	}

	if err != nil {
		t.Fatalf("deleteWorkspace failed: %v", err)
	}

	t.Logf("deleted workspace id=%s", resp.Id)
}

//Users Test

func TestNewUser_Success(t *testing.T) {	
	cfg := config.Get()
	client := NewClient(cfg.Endpoint)
	priv := keys.GenerateKey(t)
	resp, err := users.CreateUser(client, priv)
	if err != nil {
		t.Fatalf("createUser failed: %v", err)
	}
	t.Logf("created user id=%s", resp.Id)

	rId, err := users.DeleteUser(client, priv)

	if rId != resp.Id {
		t.Fatalf("deleteUser failed: %v", err)
	}

	if err != nil {
		t.Fatalf("deleteUser failed: %v", err)
	}
	t.Logf("deleted user id=%s", resp.Id)
}


func TestDeleteUserByPassword(t *testing.T) {	
	cfg := config.Get()
	client := NewClient(cfg.Endpoint)
	priv := keys.GenerateKey(t)
	resp, err := users.CreateUser(client, priv)
	if err != nil {
		t.Fatalf("createUser failed: %v", err)
	}
	t.Logf("created user id=%s", resp.Id)

	rId, err := users.DeleteUserByPassword(client, resp.Id, resp.Password)

	if rId != resp.Id {
		t.Fatalf("deleteUser failed: %v", err)
	}

	if err != nil {
		t.Fatalf("deleteUser failed: %v", err)
	}

	t.Logf("deleted user id=%s", resp.Id)
}
package test

import (	
	"testing"
	"novakeyclient"
	"github.com/core-regulus/novakey-types-go"
)

//Workspace tests
func TestSetWorkspace_Success(t *testing.T) {		
	client := novakeyclient.NewClient()
	_ , priv, _ := CreateUser(t, client, "", []novakeytypes.Workspace{}, []novakeytypes.Project{})
	defer func() {
		if _, err := DeleteUser(client, priv); err != nil {
			t.Logf("failed to delete user1 in cleanup: %v", err)
		}
	}()
	workspaceResp, err := CreateWorkspace(client, priv)
	if err != nil {
		t.Fatalf("createWorkspace failed: %v", err)
	}
	t.Logf("created workspace id=%s", workspaceResp.Id)

	projectResp, err := CreateProject(
		client, workspaceResp.Id, 
		priv, 
		[]novakeytypes.Key{
			{ Key: "TestKey", Value: "TestValue" },
		},
	)
	if err != nil {
		t.Fatalf("create project failed: %v", err)
	}
	t.Logf("created project id=%s", projectResp.Id)

	_ , priv2, _ := CreateUser(t, client, priv,
		[]novakeytypes.Workspace{
    	{Id: workspaceResp.Id, RoleCode: "root.workspace.admin" },
		},
		[]novakeytypes.Project{
    	{Id: projectResp.Id, RoleCode: "root.workspace.project.admin" },
		})

	

	defer func() {
		pId, err := DeleteProject(client, projectResp.Id, priv)

		if pId != projectResp.Id {
			t.Fatalf("delete project failed: %v", err)
		}

		if _, err := DeleteUser(client, priv2); err != nil {
			t.Logf("failed to delete user1 in cleanup: %v", err)
		}

		rId, err := DeleteWorkspace(client, workspaceResp.Id, priv)

		if rId != workspaceResp.Id {
			t.Fatalf("deleteWorkspace failed: %v", err)
		}
	}()
	
	if err != nil {
		t.Fatalf("deleteWorkspace failed: %v", err)
	}
	t.Logf("deleted workspace id=%s", workspaceResp.Id)
}


func TestDeleteWorkspaceByPassword(t *testing.T) {	
	client := novakeyclient.NewClient()
	userResp , priv, _ := CreateUser(t, client, "", []novakeytypes.Workspace{}, []novakeytypes.Project{})
	
	defer func() {
		if _, err := DeleteUser(client, priv); err != nil {
			t.Logf("failed to delete user in cleanup: %v", err)
		}
	}()

	resp, err := CreateWorkspace(client, priv)
	if err != nil {
		t.Fatalf("createWorkspace failed: %v", err)
	}
	t.Logf("created workspace id=%s", resp.Id)

	rId, err := DeleteWorkspaceByPassword(client, resp.Id, userResp.Id, userResp.Password)

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
	client := novakeyclient.NewClient()	
	resp, priv, err := CreateUser(t, client, "", []novakeytypes.Workspace{}, []novakeytypes.Project{})
	if err != nil {
		t.Fatalf("createUser failed: %v", err)
	}
	t.Logf("created user id=%s", resp.Id)
	
	rId, err := DeleteUser(client, priv)

	if rId != resp.Id {
		t.Fatalf("deleteUser failed: %v", err)
	}

	if err != nil {
		t.Fatalf("deleteUser failed: %v", err)
	}
	t.Logf("deleted user id=%s", resp.Id)
}


func TestDeleteUserByPassword(t *testing.T) {		
	client := novakeyclient.NewClient()	
	resp, _, err := CreateUser(t, client, "", []novakeytypes.Workspace{}, []novakeytypes.Project{})
	if err != nil {
		t.Fatalf("createUser failed: %v", err)
	}
	t.Logf("created user id=%s", resp.Id)

	rId, err := DeleteUserByPassword(client, resp.Id, resp.Password)

	if rId != resp.Id {
		t.Fatalf("deleteUser failed: %v", err)
	}

	if err != nil {
		t.Fatalf("deleteUser failed: %v", err)
	}

	t.Logf("deleted user id=%s", resp.Id)
}
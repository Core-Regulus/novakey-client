package test

import (	
	"testing"
	"novakeyclient"
	"github.com/core-regulus/novakey-types-go"
)

/*var launchCfg = novakeyclient.LaunchConfig{
	Backend: novakeyclient.BakendConfig{
		Endpoint: "http://localhost:5000",
	},
}*/

var launchCfg = novakeyclient.LaunchConfig{
	Backend: novakeyclient.BakendConfig{
		Endpoint: "https://novakey-api.core-regulus.com",
	},
}


//Workspace tests
func TestSetWorkspace_Success(t *testing.T) {		
	client := novakeyclient.NewClientFromLaunchConfig(launchCfg)
	_ , priv, _ := CreateUser(t, client, "", []novakeytypes.Workspace{}, []novakeytypes.Project{}, "")
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

	keysResp, err := GetProject(
		client, projectResp.Id, 
		priv,
	)
	if err != nil {
		t.Fatalf("User1 getKeys failed: %v", err)
	}
	
	tKey := keysResp.Keys[0]
	if (tKey.Key != "TestKey") || (tKey.Value != "TestValue") {
		t.Fatalf("Key error received %s - %s", tKey.Key, tKey.Value)
	}
		
	_ , priv2, _ := CreateUser(t, client, priv,
		[]novakeytypes.Workspace{
    	{Id: workspaceResp.Id, RoleCode: "root.workspace.admin" },
		},
		[]novakeytypes.Project{
		{
				Id: projectResp.Id, 
				RoleCode: "root.workspace.project.admin",
			},
		},
		"testuser1@test.com",	
	)

	keysResp2, err := GetProject(
		client, projectResp.Id, 
		priv2,
	)
	if err != nil {
		t.Fatalf("User2 getKeys failed: %v", err)
	}
	
	tKey2 := keysResp2.Keys[0]
	if (tKey2.Key != "TestKey") || (tKey2.Value != "TestValue") {
		t.Fatalf("Key error received %s - %s", tKey2.Key, tKey2.Value)
	}

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
		if _, err := DeleteUser(client, priv); err != nil {
			t.Logf("failed to delete user1 in cleanup: %v", err)
		}
	}()
	
	if err != nil {
		t.Fatalf("deleteWorkspace failed: %v", err)
	}
	t.Logf("deleted workspace id=%s", workspaceResp.Id)
}


func TestDeleteWorkspaceByPassword(t *testing.T) {	
	client := novakeyclient.NewClientFromLaunchConfig(launchCfg)
	userResp , priv, _ := CreateUser(t, client, "", []novakeytypes.Workspace{}, []novakeytypes.Project{}, "")
	
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
	client := novakeyclient.NewClientFromLaunchConfig(launchCfg)
	resp, priv, err := CreateUser(t, client, "", []novakeytypes.Workspace{}, []novakeytypes.Project{}, "")
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
	client := novakeyclient.NewClientFromLaunchConfig(launchCfg)
	resp, _, err := CreateUser(t, client, "", []novakeytypes.Workspace{}, []novakeytypes.Project{}, "")
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
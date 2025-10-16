package novakeyclient

import (
	"context"
	"log"
	"net/http"
	"strings"
	"fmt"
	novakeytypes "github.com/core-regulus/novakey-types-go"	
)

func applyConfig(client *Client, privateKey string, cfg LaunchConfig) (LaunchConfig, error) {
	res := cfg;
	userRes :=client.SetUser(context.Background(), privateKey, "", novakeytypes.SetUserRequest{
		Email: cfg.Signer.Email,
	})

	if (userRes.Error.Status != 200) {
		return res, fmt.Errorf("set user: %s", novakeytypes.FormatErrorResponse(userRes.Error))
	}
			
	if (cfg.Workspace.Name == "") {
		return res, nil
	}

	req := novakeytypes.SetWorkspaceRequest{
		Id: cfg.Workspace.Id,
		Name: cfg.Workspace.Name,		
	}
	workspaceResp := client.SetWorkspace(context.Background(), privateKey, req)
	
	if (workspaceResp.Error.Error != "") {
		return res, fmt.Errorf("set workspace: %s", novakeytypes.FormatErrorResponse(workspaceResp.Error))
	}

	res.Workspace.Id = workspaceResp.Id;
	res.Workspace.RoleCode = workspaceResp.RoleCode

	if (cfg.Workspace.Project.Name == "") {
		return res, nil
	}

	projectResp := client.SetProject(context.Background(), privateKey, novakeytypes.SetProjectRequest{
		Id: cfg.Workspace.Project.Id,
		Name: cfg.Workspace.Project.Name,
		Description: cfg.Workspace.Project.Description,
		WorkspaceId: workspaceResp.Id,
		Keys: cfg.Workspace.Project.Keys,
	})
	
	if (projectResp.Error.Error != "") {
		return res, fmt.Errorf("set project: %s", novakeytypes.FormatErrorResponse(projectResp.Error))
	}

	res.Workspace.Project.WorkspaceId = workspaceResp.Id;
	res.Workspace.Project.RoleCode = projectResp.RoleCode;
	res.Workspace.Project.Id = projectResp.Id;
	return res, nil	

}

func NewClient(cfg InitConfig) (*LaunchConfig, *Client, error) {
	launchCfg, err := Load(cfg)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
		
	client := NewClientFromLaunchConfig(*launchCfg)
	res, err := applyConfig(client, launchCfg.Signer.PrivateKey, *launchCfg)
	if err != nil {
		return launchCfg, client, err
	}
	
	err = saveLaunchFile(cfg, &res)
	return launchCfg, client, err
}


func NewClientFromLaunchConfig(launchConfig LaunchConfig) *Client {
	return &Client{
		BaseURL:    strings.TrimRight(launchConfig.Backend.Endpoint, "/"),
		HTTPClient: http.DefaultClient,
	}
}
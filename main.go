package novakeyclient

import (
	"context"
	"errors"
	"fmt"	
	"net/http"
	"strings"

	novakeytypes "github.com/core-regulus/novakey-types-go"
)

func applyConfig(privateKey string, cfg *LaunchConfig) (*LaunchConfig, error) {
	res := cfg;
			
	userRes := res.Client.SetUser(context.Background(), privateKey, "", novakeytypes.SetUserRequest{
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
	workspaceResp := res.Client.SetWorkspace(context.Background(), privateKey, req)
	
	if (workspaceResp.Error.Error != "") {
		return res, fmt.Errorf("set workspace: %s", novakeytypes.FormatErrorResponse(workspaceResp.Error))
	}

	res.Workspace.Id = workspaceResp.Id;
	res.Workspace.RoleCode = workspaceResp.RoleCode

	if (cfg.Workspace.Project.Name == "") {
		return res, nil
	}

	projectResp := res.Client.SetProject(context.Background(), privateKey, novakeytypes.SetProjectRequest{
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

func loadFromInit(cfg InitConfig) (*LaunchConfig, error) {
	launchCfg, err := LoadFromInitConfig(cfg)
	if err != nil {
		return nil, err;
	}
			
	res, err := applyConfig(launchCfg.Signer.PrivateKey, launchCfg)
	if err != nil {
		return res, err
	}
	
	err = saveLockFile(cfg, res)
	return res, err
}

func NewClient(cfg InitConfig) (*LaunchConfig, error) {	
	LaunchConfig, err := loadFromInit(cfg)
	if (err == nil) {
		return LaunchConfig, nil
	}

	if !errors.Is(err, ErrInitFileNotFound) {
		return nil, err;
	}

	launchCfg, err := LoadFromLockFile(cfg)	
	if (err == nil) {
		return launchCfg, nil
	}
	
	return launchCfg, err;
}


func NewClientFromLaunchConfig(launchConfig LaunchConfig) *Client {
	return &Client{
		BaseURL:    strings.TrimRight(launchConfig.Backend.Endpoint, "/"),
		HTTPClient: http.DefaultClient,
	}
}
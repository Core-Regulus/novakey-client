package novakeyclient

import (
	novakeytypes "github.com/core-regulus/novakey-types-go"
	"github.com/google/uuid"
)

type InitConfig struct {
	Directory string
}

type LaunchConfig struct {
	Client		*Client									`yaml:"-"`
	Backend 	BakendConfig   					`yaml:"backend"`
	Workspace novakeytypes.Workspace 	`yaml:"workspace"`	
	Signer  	novakeytypes.Signer			`yaml:"-"`
}

type UserConfig struct {
	Email 					string	`yaml:"email"`	
	PrivateKeyFile 	string  `yaml:"privateKeyFile"`
}

type BakendConfig struct {
	Endpoint    string  		`yaml:"endpoint"`	
}

type LockConfig struct {
	WorkspaceId uuid.UUID `yaml:"workspaceId"`
	ProjectId   uuid.UUID `yaml:"projectId"`
}


package test

import (
	"context"
	"github.com/core-regulus/novakey-client"
	"os"
	"path/filepath"
	"testing"

	novakeytypes "github.com/core-regulus/novakey-types-go"
)

var dir = "./testConfigs"
var key = ".privateKey.pem"

func checkFile(t *testing.T, path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("expected file %s to exist", path)
	}
}

func Test(t *testing.T) {	
	keyFilepath := filepath.Join(dir, key)
	GeneratekeyToFile(keyFilepath)
	
	launchCfg, err := novakeyclient.NewClient(novakeyclient.InitConfig{ Directory:  dir })	
	if (err != nil) {
		t.Fatalf("createClient from novakey-init failed: %v", err)
	}

	checkFile(t, keyFilepath)
	checkFile(t, filepath.Join(dir, ".novakey-user.yaml"))
	checkFile(t, filepath.Join(dir, "novakey-launch.yaml"))
		
	launchCfg, err = novakeyclient.NewClient(novakeyclient.InitConfig{ Directory:  dir })	
	if (err != nil) {
		t.Fatalf("createClient from novake-launch failed: %v", err)
	}

	defer func() {
		launchCfg.Client.DeleteUser(context.Background(), launchCfg.Signer.PrivateKey, novakeytypes.DeleteUserRequest{
				AuthEntity: novakeytypes.AuthEntity{
				Id: launchCfg.Signer.Id,
    	},	
		})
		os.Remove(keyFilepath)		
		os.Remove(filepath.Join(dir, "novakey-launch.yaml"))
	}()

		
}
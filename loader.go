package novakeyclient

import (
	"errors"
	"os"
	"path/filepath"
	novakeytypes "github.com/core-regulus/novakey-types-go"
	"gopkg.in/yaml.v3"
)

var initFilename = ".novakey-init.yaml"
var lockFilename = "novakey-lock.yaml"
var userFilename = ".novakey-user.yaml"

func stat(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	}
	return false
}

func getInitFile(cfg InitConfig) (string, error) {	
	initFile := filepath.Join(cfg.Directory,initFilename)	
	if stat(initFile) {
		return initFile, nil
	}
	
	return "", errors.New("no .novakey-init.yaml file found")	
}

func getUserFile(cfg InitConfig) (string, error) {
	userFile :=  filepath.Join(cfg.Directory, userFilename)
	
	if stat(userFile) {
		return userFile, nil
	} 
	return "", errors.New("no .novakey-user.yaml file found")	
}

func getLockFile(cfg InitConfig) (string, error) {
	lockFile :=  filepath.Join(cfg.Directory, lockFilename)
	
	if stat(lockFile) {
		return lockFile, nil
	} 
	return "", errors.New("no novakey-lock.yaml file found")	
}

func saveLockFile(cfg InitConfig, launchCfg *LaunchConfig) error {
	lockFile := filepath.Join(cfg.Directory, lockFilename)
	lockCfg := LockConfig{
		WorkspaceId: launchCfg.Workspace.Id,
		ProjectId:   launchCfg.Workspace.Project.Id,
	}
		
	data, err := yaml.Marshal(lockCfg)
  if err != nil {
    return errors.New("error in marshal lockConfig")	
  }
    
  err = os.WriteFile(lockFile, data, 0644)
  if err != nil {
    return errors.New("error writing novakey-lock.yaml file")	
  }

	return nil
}

func Load(cfg InitConfig) (*LaunchConfig, error) {	
	initFile, err := getInitFile(cfg)
	if err != nil {
		return nil, err
	}

	userFile, err := getUserFile(cfg)
	if err != nil {
		return nil, err
	}
	
	user, err := loadUserFromYaml(userFile)	
	if err != nil {
		return nil, err
	}

	res, err := loadYAML(initFile);	
	if err == nil {
		res.Signer = *user
	}
	lock, _ := loadLockYaml(cfg)
	res.Workspace.Id = lock.WorkspaceId
	res.Workspace.Project.Id = lock.ProjectId	
	res.Client = NewClientFromLaunchConfig(*res)
	return res, err;
}

func loadYAML(filename string) (*LaunchConfig, error) {	
	res := &LaunchConfig{}
	file, err := os.Open(filename)
	if err != nil {
		return nil, err;
	}
	defer file.Close()
	
	decoder := yaml.NewDecoder(file)
	
	if err := decoder.Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func loadLockYaml(cfg InitConfig) (*LockConfig, error) {		
	res := &LockConfig{}
	filename, err := getLockFile(cfg)
	if err != nil {
		return res, err
	}
	file, err := os.Open(filename)
	if err != nil {
		return res, err;
	}
	defer file.Close()
	
	decoder := yaml.NewDecoder(file)
	
	if err := decoder.Decode(&res); err != nil {
		return res, err
	}

	return res, nil
}

func loadUserFromYaml(filename string) (*novakeytypes.Signer, error) {		
	file, err := os.Open(filename)
	if err != nil {
		return nil, err;
	}
	defer file.Close()
	
	decoder := yaml.NewDecoder(file)
	cfg := UserConfig{}
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	privBytes, err := os.ReadFile(cfg.PrivateKeyFile)
	if err != nil {
		return nil, err
	}

	res := novakeytypes.Signer{
		Email: cfg.Email,
		PrivateKey: string(privBytes),
	}

	return &res, nil
}
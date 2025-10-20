package novakeyclient

import (
	"errors"
	"os"
	"path/filepath"

	novakeytypes "github.com/core-regulus/novakey-types-go"
	"gopkg.in/yaml.v3"
)


func stat(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	}
	return false
}

func getLaunchFile(cfg InitConfig) (string, error) {
	initFilename :=  filepath.Join(cfg.Directory, "novakey-init.yaml")
	launchFilename := filepath.Join(cfg.Directory, "novakey-launch.yaml")

	if stat(launchFilename) {
		return launchFilename, nil
	} 
	if stat(initFilename) {
		return initFilename, nil
	}
	return "", errors.New("no config file found")	
}

func saveLaunchFile(cfg InitConfig, launchCfg *LaunchConfig) error {
	launchFilename := filepath.Join(cfg.Directory, "novakey-launch.yaml")
	data, err := yaml.Marshal(launchCfg)
  if err != nil {
    return errors.New("error in marshal config")	
  }
    
  err = os.WriteFile(launchFilename, data, 0644)
  if err != nil {
    return errors.New("error writing config file")	
  }

	return nil
}

func getUserFile(cfg InitConfig) (string, error) {
	userFilename :=  filepath.Join(cfg.Directory, ".novakey-user.yaml")
	
	if stat(userFilename) {
		return userFilename, nil
	} 
	return "", errors.New("no user file found")	
}

func Load(cfg InitConfig) (*LaunchConfig, error) {	
	launchFilename, err := getLaunchFile(cfg)
	if err != nil {
		return nil, err
	}
	userFilename, err := getUserFile(cfg)
	if err != nil {
		return nil, err
	}
	
	user, err := loadUserFromYaml(userFilename)	
	if err != nil {
		return nil, err
	}

	res, err := loadYAML(launchFilename);	
	if err == nil {
		res.Signer = *user
	}
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
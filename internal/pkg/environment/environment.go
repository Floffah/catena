package environment

import (
	"os"
	"os/exec"
)

type Environment struct {
	Config Config
	GitBin string
}

func LoadEnvironment() (*Environment, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	// make sure the git root directory exists
	_, err = os.Stat(cfg.CatenaGitRoot)
	if os.IsNotExist(err) {
		err = os.MkdirAll(cfg.CatenaGitRoot, 0755)
		if err != nil {
			return nil, err
		}
	}

	gitBin, err := exec.LookPath("git")
	if err != nil {
		return nil, err
	}

	return &Environment{
		Config: *cfg,
		GitBin: gitBin,
	}, nil
}

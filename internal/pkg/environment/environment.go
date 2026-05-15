package environment

import (
	"os"
	"os/exec"
)

type Environment struct {
	Config Config
	GitBin string
}

func LoadEnvironment() (Environment, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return Environment{}, err
	}

	// make sure the git root directory exists
	_, err = os.Stat(cfg.CatenaGitRoot)
	if os.IsNotExist(err) {
		err = os.MkdirAll(cfg.CatenaGitRoot, 0750)
		if err != nil {
			return Environment{}, err
		}
	}

	gitBin, err := exec.LookPath("git")
	if err != nil {
		return Environment{}, err
	}

	return Environment{
		Config: cfg,
		GitBin: gitBin,
	}, nil
}

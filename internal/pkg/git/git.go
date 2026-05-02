package git

import "os/exec"

// Git client, currently just a backend for gitstore and git binary, but eventually will incorporate go-git
type Git struct {
	BinaryPath string
}

func NewGit(binaryPath string) Git {
	return Git{
		BinaryPath: binaryPath,
	}
}

func (g Git) Init(repoPath string) error {
	err := exec.Command(g.BinaryPath, "init", repoPath).Run()
	if err != nil {
		return err
	}

	return nil
}

func (g Git) InitBare(repoPath string, defaultBranch string) error {
	err := exec.Command(g.BinaryPath, "init", "--bare", "--initial-branch="+defaultBranch, repoPath).Run()
	if err != nil {
		return err
	}

	return nil
}

func (g Git) Path() string {
	return g.BinaryPath
}

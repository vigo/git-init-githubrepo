package command

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

// sentinel errors.
var (
	ErrCWDIsEmpty        = errors.New("current working directory is empty")
	ErrAlreadyInAGitRepo = errors.New("you are already in a git repo")
)

func (k *cmd) setCWD() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get current working directory: %w", err)
	}
	k.cwd = cwd

	return nil
}

func (k *cmd) checkDefaults() error {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		return fmt.Errorf("git is not installed: %w", err)
	}
	k.gitPath = gitPath

	if err = k.setCWD(); err != nil {
		return fmt.Errorf("could not set current working directory: %w", err)
	}

	if k.cwd == "" {
		return ErrCWDIsEmpty
	}

	userFullName, err := k.runGITCommand("config", "user.name")
	if err == nil {
		k.gitUserFullName = userFullName
	}

	userEmail, err := k.runGITCommand("config", "user.email")
	if err == nil {
		k.gitUserEmail = userEmail
	}

	gitHubUserName, err := k.runGITCommand("config", "github.user")
	if err == nil {
		k.gitHubUserName = gitHubUserName
	}

	return nil
}

func (k *cmd) runGITCommand(args ...string) (string, error) {
	var out strings.Builder

	joinedArgs := slices.Concat([]string{"git"}, args)

	execCmd := &exec.Cmd{
		Path:   k.gitPath,
		Args:   joinedArgs,
		Stdout: &out,
		Stderr: &out,
	}

	if err := execCmd.Start(); err != nil {
		return "", fmt.Errorf("can not start git command: %w", err)
	}

	if err := execCmd.Wait(); err != nil {
		return "", fmt.Errorf("can npt wait git command: %w", err)
	}

	return strings.TrimSpace(out.String()), nil
}

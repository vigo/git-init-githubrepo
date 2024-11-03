package main

import (
	"fmt"
	"os"

	"github.com/vigo/git-init-githubrepo/internal/command"
)

func main() {
	cmd, err := command.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if err = cmd.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

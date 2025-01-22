package main

import (
	"fmt"
	"os"

	"github.com/gogodjzhu/listen-tube/internal/pkg/cmd"
	"github.com/gogodjzhu/listen-tube/internal/pkg/cmd/root"
	_ "github.com/gogodjzhu/listen-tube/internal/pkg/log"
)

type exitCode int

const (
	exitOK    exitCode = 0
	exitError exitCode = 1
)

func main() {
	code := mainRun()
	os.Exit(int(code))
}

func mainRun() exitCode {
	cmdFactory := cmd.NewFactory()

	mainCmd, err := root.NewCmdRoot(cmdFactory)
	if err != nil {
		fmt.Fprintln(cmdFactory.IOStreams.Out, err)
		return exitError
	}
	if _, err := mainCmd.ExecuteC(); err != nil {
		return exitError
	}
	return exitOK
}
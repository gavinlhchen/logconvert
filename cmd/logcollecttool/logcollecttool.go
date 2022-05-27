package main

import (
	"logconvert/internal/logcollecttool/cmd"
	"os"
)

func main() {
	command := cmd.NewDefaultLogCollectToolCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

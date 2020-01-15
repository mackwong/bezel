package main

import (
	"gitlab.bj.sensetime.com/diamond/bezel/cmd"
	"os"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

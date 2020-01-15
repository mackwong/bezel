package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var debug bool

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:               "bezel",
		Short:             "Bezel control interface.",
		SilenceUsage:      true,
		DisableAutoGenTag: true,
		Long:              "",
		PersistentPreRun:  bezelPersistentPreRun,
	}
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "show debug information")
	rootCmd.AddCommand(NewCreateCmd(), NewParseCmd())

	return rootCmd
}

func bezelPersistentPreRun(cmd *cobra.Command, args []string) {
	log.SetFormatter(&log.TextFormatter{DisableTimestamp: true})
	log.SetLevel(log.InfoLevel)
	if debug {
		log.SetLevel(log.DebugLevel)
	}
}

package cmd

import (
	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "bezel",
		Short: "Bezel control interface.",
		SilenceUsage:      true,
		DisableAutoGenTag: true,
		Long: "",
		PersistentPreRun: bezelPersistentPreRun,
	}
	rootCmd.PersistentFlags().Bool("debug", false, "show debug information")

	rootCmd.AddCommand(NewCreateCmd())
	//rootCmd.AddCommand(NewDeleteCmd)

	return rootCmd
}

func bezelPersistentPreRun(cmd *cobra.Command, args []string) {
	log.SetFormatter(&log.TextFormatter{DisableTimestamp:true})
	log.SetLevel(log.InfoLevel)
}

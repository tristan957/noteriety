package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.PersistentFlags().StringP("repository", "C", "", "path to a Noteriety repository")
	viper.BindPFlag("repository.path", rootCmd.PersistentFlags().Lookup("repository"))
}

var rootCmd = &cobra.Command{
	Use:   "notoriety",
	Short: "Notoriety is a note creation and management tool",
	Long:  "A note creation and management tool based around relationships without the need for closed source software services where you don't own your data.",
}

// Execute the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

package cmd

import (
	"github.com/chaostreff-flensburg/moc/config"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var executeMigrate = false
var executeSeed = false

// rootCmd will run the log streamer
var rootCmd = cobra.Command{
	Use:  "moc",
	Long: "A service that will serve a restufl massage operation center",
	Run: func(cmd *cobra.Command, args []string) {
		execWithConfig(cmd, serve)
	},
}

// RootCmd will add flags and subcommands to the different commands
func RootCmd() *cobra.Command {
	rootCmd.PersistentFlags().BoolVarP(&executeMigrate, "migrate", "m", false, "migrate database")
	rootCmd.PersistentFlags().BoolVarP(&executeSeed, "seed", "s", false, "seed database")
	rootCmd.AddCommand(&serveCmd)
	return &rootCmd
}

// execWithConfig load config from env
func execWithConfig(cmd *cobra.Command, fn func(config *config.Config)) {
	logrus.Info("Read Config...")
	config := config.ReadConfig()

	fn(config)
}

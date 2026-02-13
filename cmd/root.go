package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "lume-fleet",
	Short: "Manage a fleet of Lume VMs declaratively",
	Long:  "lume-fleet reads a fleet.yml file and manages multiple Lume VMs as a group.",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "fleet.yml", "path to fleet config file")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

package main

import (
	"errors"
	"go_chain/config"
	"go_chain/server"

	"github.com/spf13/cobra"
)

var run = &cobra.Command{
	Use:   "run",
	Short: "run a mining with a provided block info",
	RunE: func(cmd *cobra.Command, args []string) error {
		server := server.New(&config.Config)
		if err := server.Start(); err != nil {
			return errors.New("Failed to start the blockchain node")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(run)
}

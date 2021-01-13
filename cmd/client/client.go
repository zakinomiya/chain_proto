package main

import (
	"chain_proto/gateway"
	"chain_proto/peer"

	"github.com/spf13/cobra"
)

var run = &cobra.Command{
	Use:   "block",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := gateway.NewClient()
		c.AddNeighbour(peer.New("", "tcp"))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(run)
}

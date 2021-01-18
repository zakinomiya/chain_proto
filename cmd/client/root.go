package main

import (
	"chain_proto/gateway"
	"chain_proto/peer"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	server string
	c      *gateway.Client
)

var rootCmd = &cobra.Command{
	Use:   "client",
	Short: "client command provides commands for a blockchain server",
	Long:  ``,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cmd.PersistentFlags().StringVarP(&server, "server", "s", "localhost:9000", "Request endpoint with the port specified.")
		c = gateway.NewClient()
		c.AddNeighbour(peer.New(server, "tcp"))

		fmt.Printf("Requesting to %s\n", server)
		return nil
	},
}

func init() {
	initBlockCmd()

	rootCmd.AddCommand(blockCmds)
}

func execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

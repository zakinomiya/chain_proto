package main

import (
	"chain_proto/config"
	"chain_proto/server"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var run = &cobra.Command{
	Use:   "run",
	Short: "run a mining with a provided block info",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := server.New(config.Config)
		if err != nil {
			log.Println("info: Failed to initialise the server")
			return err
		}

		if err := server.Start(); err != nil {
			log.Println("Failed to start the blockchain node")
			return err
		}

		waitExit()

		server.Stop()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(run)
}

func waitExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	i := <-c
	log.Printf("info: received signal %s. Stopping", i)
}

package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"go_chain/wallet"

	"github.com/spf13/cobra"
)

var run = &cobra.Command{
	Use:   "run",
	Short: "run a mining with a provided block info",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("no args found")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		privKeyStr := args[0]

		privKeyBytes, err := hex.DecodeString(privKeyStr)
		if err != nil {
			return errors.New("Invalid form of private key")
		}

		w, err := wallet.RestoreWallet(privKeyBytes)
		if err != nil {
			return errors.New("Failed to inistialise wallet")
		}

		fmt.Printf("private key: %s\npublic key: %s\n", w.PrivKeyStr(), w.PubKeyStr())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(run)
}

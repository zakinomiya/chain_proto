package wallet

import (
	"encoding/hex"
	"errors"
	"fmt"
	"go_chain/wallet"

	"github.com/spf13/cobra"
)

var newKeyPair = &cobra.Command{
	Use:   "keypair",
	Short: "create a new keypair. The created keypair is not stored but just output to stdout.",
	RunE: func(cmd *cobra.Command, args []string) error {
		w, err := wallet.New()
		if err != nil {
			return errors.New("Failed to inistialise wallet")
		}

		fmt.Printf("new private key: %s\nnew public key: %s\n", w.PrivKeyStr(), w.PubKeyStr())
		return nil
	},
}

var restoreKeyPair = &cobra.Command{
	Use:   "restore",
	Short: "restore a new keypair from private key",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("First argument must be private key")
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
	rootCmd.AddCommand(newKeyPair)
	rootCmd.AddCommand(restoreKeyPair)
}

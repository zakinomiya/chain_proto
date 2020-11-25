package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"go_chain/wallet"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const (
	keypairDirName     = "keypairs"
	defaultKeyPairName = "keypair"
)

var keypairName string

var newKeyPair = &cobra.Command{
	Use:   "keypair",
	Short: "Create a new keypair.",
	RunE: func(cmd *cobra.Command, args []string) error {
		w, err := wallet.New()
		if err != nil {
			return errors.New("Failed to inistialise wallet")
		}

		shouldSave, err := cmd.Flags().GetBool("save")
		if err != nil {
			return nil
		}

		if shouldSave {
			return storeKeypair(keypairDirName, keypairName, w)
		}

		fmt.Printf("private key=%s, public key=%s\n", w.PrivKeyStr(), w.PubKeyStr())
		return nil
	},
}

var restoreKeyPair = &cobra.Command{
	Use:   "restore",
	Short: "Restore a keypair from private key",
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

		shouldSave, err := cmd.Flags().GetBool("save")
		if err != nil {
			return nil
		}

		if shouldSave {
			return storeKeypair(keypairDirName, keypairName, w)
		}

		fmt.Printf("private key=%s, public key=%s\n", w.PrivKeyStr(), w.PubKeyStr())
		return nil
	},
}

func init() {
	setCommonFlags(newKeyPair, restoreKeyPair)

	rootCmd.AddCommand(newKeyPair, restoreKeyPair)
}

func storeKeypair(dir string, keypairName string, wallet *wallet.Wallet) error {
	keypairYml := []yaml.MapItem{
		{
			Key:   "private_key",
			Value: wallet.PrivKeyStr(),
		},
		{
			Key:   "public_key",
			Value: wallet.PubKeyStr(),
		},
	}

	ymlByte, err := yaml.Marshal(keypairYml)
	if err != nil {
		return err
	}

	if _, err := os.Stat(filepath.Join(keypairDirName, keypairName+".yml")); !os.IsNotExist(err) {
		return errors.New("A keypair file with the same name already exists")
	}

	if _, err := os.Stat(keypairDirName); os.IsNotExist(err) {
		os.Mkdir(keypairDirName, 0755)
	}

	if err := ioutil.WriteFile(filepath.Join(dir, keypairName+".yml"), ymlByte, 0400); err != nil {
		fmt.Println("JJJJ")
		return err
	}

	fmt.Printf("A new keypair file (%s.yml) has been successfully created\n", keypairName)
	return nil
}

func setCommonFlags(cmds ...*cobra.Command) {
	for _, cmd := range cmds {
		cmd.Flags().BoolP("save", "s", false, "Keypair will be saved to a file")
		cmd.Flags().StringVarP(&keypairName, "name", "n", "", "Filename for a new keypair")
	}

	if keypairName == "" {
		keypairName = defaultKeyPairName
	}
}

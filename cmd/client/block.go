package main

import (
	"chain_proto/common"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var blockCmds = &cobra.Command{
	Use: "block",
}

var byHashCmd = &cobra.Command{
	Use: "get",
	RunE: func(cmd *cobra.Command, args []string) error {
		hash := args[0]
		if hash == "" {
			fmt.Println("ERROR: please specify block hash")
			return errors.New("")
		}

		hashStr, err := hex.DecodeString(hash)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		b, err := c.GetBlockByHash(ctx, common.ReadByteInto32(hashStr))
		if err != nil {
			return err
		}

		j, err := json.Marshal(b)
		if err != nil {
			return err
		}

		fmt.Sprintln("Found block")
		fmt.Printf("%s\n", j)
		return nil
	},
}

func initBlockCmd() {
	blockCmds.AddCommand(byHashCmd)
}

package main

import (
	"chain_proto/common"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var txCmds = &cobra.Command{
	Use: "tx",
}

var byTxHashCmd = &cobra.Command{
	Use: "get",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Please provide transaction hash")
		}

		hash := strings.Split(args[0], "=")[1]

		hByte, err := hex.DecodeString(hash)
		if err != nil {
			return fmt.Errorf("Invalid hash")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		tx, err := c.GetTransactionByHash(ctx, common.ReadByteInto32(hByte))
		if err != nil {
			return fmt.Errorf("%s", err)
		}

		j, err := json.Marshal(tx)
		if err != nil {
			return fmt.Errorf("failed to get transaction.")
		}

		fmt.Printf(string(j))
		return nil
	},
}

func initTxCmd() {
	txCmds.AddCommand(byTxHashCmd)
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var accCmds = &cobra.Command{
	Use: "account",
}

var getCmd = &cobra.Command{
	Use: "get",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Please provide arguments.\n %s", valExamples)
		}

		val := args[0]
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		addr := strings.Split(val, "=")[1]

		acc, err := c.GetAccount(ctx, addr)
		if err != nil {
			return fmt.Errorf("failed to get account. err=%+v\n", err)
		}

		j, err := json.Marshal(acc)
		if err != nil {
			return fmt.Errorf("failed to get account. err=%+v\n", err)
		}

		fmt.Println(string(j))
		return nil
	},
}

func initAccountCmd() {
	accCmds.AddCommand(getCmd)
}

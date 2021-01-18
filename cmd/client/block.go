package main

import (
	"chain_proto/common"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var blockCmds = &cobra.Command{
	Use: "block",
}

const valExamples = `argument examples are below
	- hash=[hash] 
		ex) hash=000009fb2a2b1ba5645f2469eb7ab3c67c5695b4aea78e90ca92d735ba486345
	- height=[height]
		ex) height=10
	- range=[start,end]  (Note: start and end values are inclusive)
		ex) range=0,5
`

var byHashCmd = &cobra.Command{
	Use: "get",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Please provide arguments.\n %s", valExamples)
		}

		val := args[0]
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var jsonStr []byte
		var err error
		switch strs := strings.Split(val, "="); {
		case strs[0] == "hash":
			jsonStr, err = getBlockByHash(ctx, strs[1])
		case strs[0] == "height":
			h, err := strconv.Atoi(strs[1])
			if err != nil {
				return fmt.Errorf("Invalid height.\n %s", valExamples)
			}
			jsonStr, err = getBlockByHeight(ctx, uint32(h))
		case strs[0] == "range":
			nums := strings.Split(strs[1], ",")
			if len(nums) != 2 {
				return fmt.Errorf("Invalid range.\n %s", valExamples)
			}
			s, err := strconv.Atoi(nums[0])
			if err != nil {
				return fmt.Errorf("Invalid start height.\n %s", valExamples)
			}
			e, err := strconv.Atoi(nums[1])
			if err != nil {
				return fmt.Errorf("Invalid end height.\n %s", valExamples)
			}
			jsonStr, err = getBlocksByRange(ctx, uint32(s), uint32(e))

		default:
			return fmt.Errorf("Invalid argument.\n %s", valExamples)
		}

		if err != nil {
			return err
		}

		fmt.Println(string(jsonStr))
		return nil
	},
}

func getBlockByHash(ctx context.Context, hash string) ([]byte, error) {
	hashByte, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}

	b, err := c.GetBlockByHash(ctx, common.ReadByteInto32(hashByte))
	if err != nil {
		return nil, err
	}

	return json.Marshal(b)
}

func getBlockByHeight(ctx context.Context, height uint32) ([]byte, error) {
	b, err := c.GetBlockByHeight(ctx, height)
	if err != nil {
		return nil, err
	}

	return json.Marshal(b)
}

func getBlocksByRange(ctx context.Context, start uint32, end uint32) ([]byte, error) {
	blks, err := c.GetBlockByRange(ctx, start, end)
	if err != nil {
		return nil, err
	}

	return json.Marshal(blks)
}

func initBlockCmd() {
	blockCmds.AddCommand(byHashCmd)
}

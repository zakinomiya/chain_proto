package main

import (
	"go_chain/common"
	"go_chain/config"
)

func init() {
	common.LogginSettings(config.Config.LogFile, config.Config.DefaultLogLevel)
}

func main() {
	execute()
}

package main

import (
	"chain_proto/common"
	"chain_proto/config"
)

func init() {
	common.LogginSettings(config.Config.LogFile, config.Config.DefaultLogLevel)
}

func main() {
	execute()
}

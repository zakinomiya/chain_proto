package config

import (
	"log"
	"os"

	"gopkg.in/ini.v1"
)

type DBConfig struct {
	Path string
}

type RpcConfig struct {
	addr string
}

type ConfigSettings struct {
	ChainID         uint32
	LogFile         string
	DefaultLogLevel string
	minerSecretKey  [32]byte
	MinerPubKey     [32]byte
	DBConf          *DBConfig
	RpcConf         *RpcConfig
}

var Config ConfigSettings

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Printf("error: Failed to read config file. ERROR: %v", err)
		os.Exit(1)
	}

	Config = ConfigSettings{
		ChainID:         uint32(cfg.Section("chain_info").Key("chain_id").InUint(1995, []uint{})),
		LogFile:         cfg.Section("general").Key("log_file").String(),
		DefaultLogLevel: cfg.Section("general").Key("default_log_level").String(),
	}
}

func (c *ConfigSettings) MinerSecretKey() [32]byte {
	return c.minerSecretKey
}

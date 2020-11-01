package config

import (
	"log"
	"os"

	"gopkg.in/ini.v1"
)

type DBConfig struct {
	Path   string
	Driver string
}

type RpcConfig struct {
	addr string
}

type ConfigSettings struct {
	ChainID         uint32
	LogFile         string
	DefaultLogLevel string
	MinerSecretKey  string
	*DBConfig
	*RpcConfig
}

var Config ConfigSettings

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Printf("error: Failed to read config file. ERROR: %v", err)
		os.Exit(1)
	}

	Config = ConfigSettings{
		// general
		LogFile:         cfg.Section("general").Key("log_file").String(),
		DefaultLogLevel: cfg.Section("general").Key("default_log_level").String(),
		// chain info
		ChainID: uint32(cfg.Section("chain_info").Key("chain_id").InUint(1995, []uint{})),
		//miner
		MinerSecretKey: cfg.Section("miner").Key("secret_key").String(),
		// db
		DBConfig: &DBConfig{
			Path:   cfg.Section("db").Key("db_path").String(),
			Driver: cfg.Section("db").Key("driver").String(),
		},
	}
}

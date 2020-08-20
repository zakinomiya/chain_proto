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
	LogFile   string
	secretKey string
	DBConf    *DBConfig
	RpcConf   *RpcConfig
}

var Config ConfigSettings

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Printf("Failed to read config file. ERROR: %v", err)
		os.Exit(1)
	}

	Config = ConfigSettings{
		LogFile: cfg.Section("general").Key("log_file").String(),
	}
}

func (c *ConfigSettings) SecretKey() string {
	return c.secretKey
}

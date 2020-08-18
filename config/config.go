package config

import (
	"log"
	"os"

	"gopkg.in/ini.v1"
)

type ConfigSettings struct {
	LogFile string
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

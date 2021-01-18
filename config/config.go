package config

import (
	"chain_proto/common"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Log struct {
	LogFile         string `yaml:"log_file"`
	DefaultLogLevel string `yaml:"default_log_level"`
}

type ChainInfo struct {
	ChainID string `yaml:"chain_id"`
}

type Miner struct {
	Enabled       bool   `yaml:"enabled"`
	Concurrent    bool   `yaml:"concurrent"`
	MaxWorkersNum int    `yaml:"max_workers_num"`
	SecretKeyStr  string `yaml:"secret_key_str"`
}

type Db struct {
	DbPath   string `yaml:"db_path"`
	Driver   string `yaml:"driver"`
	SQLPath  string `yaml:"sqlPath"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type ServerInfo struct {
	Port    string `yaml:"port"`
	Enabled bool   `yaml:"enabled"`
}

type Network struct {
	RPC       ServerInfo `yaml:"rpc"`
	HTTP      ServerInfo `yaml:"http"`
	Websocket ServerInfo `yaml:"websocket"`
}

type ConfigFile struct {
	Configurations `yaml:"configurations"`
}

type Configurations struct {
	Log       `yaml:"log"`
	ChainInfo `yaml:"chain_info"`
	Miner     `yaml:"miner"`
	Db        `yaml:"db"`
	Network   `yaml:"network"`
}

var Config *Configurations

const (
	MaxDecimalDigit = 3
)

func init() {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		log.Fatalln("Please set GOPATH to run the server")
	}

	conf, err := readConfig(gopath + "/src/chain_proto/config/config.yaml")
	if err != nil {
		log.Printf("error: Failed to read config file. ERROR: %v.\n", err)
		os.Exit(1)
	}

	if ok := common.IsValidPort(conf.RPC.Port); ok == false {
		os.Exit(1)
	}

	if ok := common.IsValidPort(conf.HTTP.Port); ok == false {
		os.Exit(1)
	}

	if ok := common.IsValidPort(conf.Websocket.Port); ok == false {
		os.Exit(1)
	}

	Config = conf
}

func readConfig(path string) (*Configurations, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	configFile := &ConfigFile{}

	if err := yaml.Unmarshal(data, configFile); err != nil {
		return nil, err
	}

	// log.Printf("debug: configuration=%+v", configFile)
	return &configFile.Configurations, nil
}

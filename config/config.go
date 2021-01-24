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
	SQLPath  string `yaml:"sql_path"`
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
	Seeds     []struct {
		Addr    string `yaml:"addr"`
		Network string `yaml:"network"`
	} `yaml:"seeds"`
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

	for _, p := range []ServerInfo{conf.RPC, conf.HTTP, conf.Websocket} {
		if !p.Enabled {
			continue
		}

		if ok := common.IsValidPort(p.Port); ok == false {
			os.Exit(1)
		}
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

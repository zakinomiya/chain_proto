package config

import (
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
	ChainID uint32 `yaml:"chain_id"`
}

type Miner struct {
	SecretKeyStr string `yaml:"secret_key_str"`
}

type Db struct {
	DbPath   string `yaml:"db_path"`
	Driver   string `yaml:"driver"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}
type Network struct {
	RPCPort       uint32 `yaml:"rpc_port"`
	HTTPPort      uint32 `yaml:"http_port"`
	WebsockerPort uint32 `yaml:"websocket_port"`
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

func init() {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		log.Println("Please set GOPATH to run the server")
	}

	conf, err := readFromYaml(gopath + "/src/go_chain/config/config.yaml")
	if err != nil {
		log.Printf("error: Failed to read config file. ERROR: %v.\n", err)
		os.Exit(1)
	}

	Config = conf
}

func readFromYaml(path string) (*Configurations, error) {
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

package config

import (
	"fmt"
	"sync"
)

const DefaultConfigTemplate = `
wallet="{{ .Wallet }}"
planetmint-rpc-host="{{ .PlanetmintRPCHost }}"
r2p-host="{{ .R2PHost }}"
cron="{{ .Cron }}"
rpc-host="{{ .RPCHost }}"
rpc-user="{{ .RPCUser }}"
rpc-pass="{{ .RPCPass }}"
`

type Config struct {
	Wallet            string `mapstructure:"wallet"`
	PlanetmintRPCHost string `mapstructure:"planetmint-rpc-host"`
	R2PHost           string `mapstructure:"r2p-host"`
	Cron              string `mapstructure:"cron"`
	RPCHost           string `mapstructure:"rpc-host"`
	RPCUser           string `mapstructure:"rpc-user"`
	RPCPass           string `mapstructure:"rpc-pass"`
}

var (
	config     *Config
	initConfig sync.Once
)

// DefaultConfig returns distribution-service default config
func DefaultConfig() *Config {
	return &Config{
		Wallet:            "dao",
		PlanetmintRPCHost: "127.0.0.1:9090",
		R2PHost:           "planetmint-go-testnet-3.rddl.io",
		Cron:              "* * * * * *",
		RPCHost:           "planetmint-go-testnet-3.rddl.io:18884",
		RPCUser:           "user",
		RPCPass:           "password",
	}
}

func GetConfig() *Config {
	initConfig.Do(func() {
		config = DefaultConfig()
	})
	return config
}

func (c *Config) GetElementsURL() string {
	url := fmt.Sprintf("http://%s:%s@%s/wallet/%s", c.RPCUser, c.RPCPass, c.RPCHost, c.Wallet)
	return url
}

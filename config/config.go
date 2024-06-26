package config

import (
	"fmt"
	"sync"

	log "github.com/rddl-network/go-logger"
)

const DefaultConfigTemplate = `
wallet="{{ .Wallet }}"
planetmint-rpc-host="{{ .PlanetmintRPCHost }}"
r2p-host="{{ .R2PHost }}"
certs-path="{{ .CertsPath }}"
cron="{{ .Cron }}"
rpc-host="{{ .RPCHost }}"
rpc-user="{{ .RPCUser }}"
rpc-pass="{{ .RPCPass }}"
shamir-host="{{ .ShamirHost }}"
confirmations={{ .Confirmations }}
fund-address="{{ .FundAddress }}"
asset="{{ .Asset }}"
log-level="{{ .LogLevel }}"
`

type Config struct {
	Wallet            string `mapstructure:"wallet"`
	PlanetmintRPCHost string `mapstructure:"planetmint-rpc-host"`
	R2PHost           string `mapstructure:"r2p-host"`
	CertsPath         string `mapstructure:"certs-path"`
	Cron              string `mapstructure:"cron"`
	RPCHost           string `mapstructure:"rpc-host"`
	RPCUser           string `mapstructure:"rpc-user"`
	RPCPass           string `mapstructure:"rpc-pass"`
	ShamirHost        string `mapstructure:"shamir-host"`
	Confirmations     int    `mapstructure:"confirmations"`
	FundAddress       string `mapstructure:"fund-address"`
	Asset             string `mapstructure:"asset"`
	LogLevel          string `mapstructure:"log-level"`
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
		R2PHost:           "https://testnet-r2p.rddl.io",
		CertsPath:         "./certs/",
		Cron:              "* * * * * *",
		RPCHost:           "localhost:18884",
		RPCUser:           "user",
		RPCPass:           "password",
		ShamirHost:        "https://localhost:9091",
		Confirmations:     10,
		FundAddress:       "",
		Asset:             "7add40beb27df701e02ee85089c5bc0021bc813823fedb5f1dcb5debda7f3da9",
		LogLevel:          log.ERROR,
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

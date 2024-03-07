package config

import "sync"

const DefaultConfigTemplate = `
wallet="{{ .Wallet }}"
planetmint-rpc-host="{{ .PlanetmintRPCHost }}"
r2p-host="{{ .R2PHost }}"
cron="{{ .Cron }}"
`

type Config struct {
	Wallet            string `mapstructure:"wallet"`
	PlanetmintRPCHost string `mapstructure:"planetmint-rpc-host"`
	R2PHost           string `mapstructure:"r2p-host"`
	Cron              string `mapstructure:"cron"`
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
	}
}

func GetConfig() *Config {
	initConfig.Do(func() {
		config = DefaultConfig()
	})
	return config
}

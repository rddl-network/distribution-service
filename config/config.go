package config

import "sync"

const DefaultConfigTemplate = `
wallet="{{ .Wallet }}"
planetmint-rpc-host="{{ .PlanetmintRPCHost }}"
`

type Config struct {
	Wallet            string `mapstructure:"wallet"`
	PlanetmintRPCHost string `mapstructure:"planetmint-rpc-host"`
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
	}
}

func GetConfig() *Config {
	initConfig.Do(func() {
		config = DefaultConfig()
	})
	return config
}

package config

import (
	"bytes"
	"log"
	"os"
	"text/template"

	"github.com/spf13/viper"
)

func LoadConfig(path string) (cfg *Config, err error) {
	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigName("app")
	v.SetConfigType("toml")

	v.AutomaticEnv()

	err = v.ReadInConfig()
	if err == nil {
		cfg = GetConfig()
		cfg.Wallet = v.GetString("wallet")
		cfg.PlanetmintRPCHost = v.GetString("planetmint-rpc-host")
		cfg.R2PHost = v.GetString("r2p-host")
		cfg.CertsPath = v.GetString("certs-path")
		cfg.Cron = v.GetString("cron")
		cfg.RPCHost = v.GetString("rpc-host")
		cfg.RPCUser = v.GetString("rpc-user")
		cfg.RPCPass = v.GetString("rpc-pass")
		cfg.ShamirHost = v.GetString("shamir-host")
		cfg.Confirmations = v.GetInt("confirmations")
		cfg.FundAddress = v.GetString("fund-address")
		cfg.Asset = v.GetString("asset")
		cfg.LogLevel = v.GetString("log-level")
		cfg.DataPath = v.GetString("data-path")
		cfg.AdvisoryCron = v.GetString("advisory-cron")
		cfg.TestnetMode = v.GetBool("testnet-mode")
		cfg.TestnetAddress = v.GetString("testnet-address")
		cfg.PlmntBlocksPerDay = v.GetInt64("plmnt_blocks_per_day")
		cfg.PlmntDistributionOffset = v.GetInt64("plmnt_distribution_offset")
		cfg.DistributionSettlementOffset = v.GetInt64("distribution_settlement_offset")
		return
	}
	log.Println("no config file found.")

	tmpl := template.New("appConfigFileTemplate")
	configTemplate, err := tmpl.Parse(DefaultConfigTemplate)
	if err != nil {
		return
	}

	var buffer bytes.Buffer
	err = configTemplate.Execute(&buffer, GetConfig())
	if err != nil {
		return
	}

	err = v.ReadConfig(&buffer)
	if err != nil {
		return
	}
	err = v.SafeWriteConfig()
	if err != nil {
		return
	}

	log.Println("default config file created. please adapt it and restart the application. exiting...")
	os.Exit(0)
	return
}

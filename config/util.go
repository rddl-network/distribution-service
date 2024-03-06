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
		cfg.PlanetmintRPCHost = viper.GetString("planetmint-rpc-host")
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

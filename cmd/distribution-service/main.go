package main

import (
	"log"

	"github.com/rddl-network/distribution-service/config"
	"github.com/rddl-network/distribution-service/service"
	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	config, err := config.LoadConfig("./")
	if err != nil {
		log.Fatalf("fatal error loading config file: %s", err)
	}

	db, err := leveldb.OpenFile("./data", nil)
	if err != nil {
		db.Close()
		log.Fatal(err)
	}
	defer db.Close()

	pmClient := service.NewPlanetmintClient(config.PlanetmintRPCHost)
	eClient := service.NewElementsClient()
	r2pClient := service.NewR2PClient(config.R2PHost)
	shamirClient := service.NewShamirClient(config.ShamirHost)
	service := service.NewDistributionService(pmClient, eClient, r2pClient, shamirClient, db)

	if err = service.Run(config.Cron); err != nil {
		log.Panicf("error occurred while spinning up service: %v", err)
	}
}

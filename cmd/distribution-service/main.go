package main

import (
	"log"

	"github.com/rddl-network/distribution-service/config"
	"github.com/rddl-network/distribution-service/service"
)

func main() {
	config, err := config.LoadConfig("./")
	if err != nil {
		log.Fatalf("fatal error loading config file: %s", err)
	}

	pmClient := service.NewPlanetmintClient(config.PlanetmintRPCHost)
	eClient := service.NewElementsClient()
	service := service.NewDistributionService(pmClient, eClient)

	if err = service.Run(); err != nil {
		log.Panicf("error occurred while spinning up service: %v", err)
	}
}

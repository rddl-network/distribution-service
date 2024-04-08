package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/rddl-network/distribution-service/config"
	"github.com/rddl-network/distribution-service/service"
	r2p "github.com/rddl-network/rddl-2-plmnt-service/client"
	shamir "github.com/rddl-network/shamir-coordinator-service/client"
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
	r2pClient := r2p.NewR2PClient(config.R2PHost, &http.Client{})
	shamirClient := shamir.NewShamirCoordinatorClient(config.ShamirHost, &http.Client{})
	service := service.NewDistributionService(pmClient, eClient, r2pClient, shamirClient, db)

	// If flag distribute=true run service.Distribute function once and exit
	distribute := flag.Bool("distribute", false, "Run Distribute function once and exit")
	flag.Parse()
	if *distribute {
		service.Distribute()
		return
	}

	if err = service.Run(config.Cron); err != nil {
		log.Panicf("error occurred while spinning up service: %v", err)
	}
}

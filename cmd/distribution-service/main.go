package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rddl-network/distribution-service/config"
	"github.com/rddl-network/distribution-service/service"
	"github.com/rddl-network/go-utils/tls"
	r2p "github.com/rddl-network/rddl-2-plmnt-service/client"
	shamir "github.com/rddl-network/shamir-coordinator-service/client"
	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	config, err := config.LoadConfig("./")
	if err != nil {
		log.Fatalf("fatal error loading config file: %s", err)
	}

	var help bool
	var distribute string
	flag.BoolVar(&help, "help", false, "Lists command line options")
	flag.StringVar(&distribute, "distribute", "", "Options: 'advisories' oder 'validators'")
	flag.Parse()

	if help {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	db, err := leveldb.OpenFile("./data", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	pmClient := service.NewPlanetmintClient(config.PlanetmintRPCHost)
	eClient := service.NewElementsClient()
	r2pClient := r2p.NewR2PClient(config.R2PHost, &http.Client{})
	mTLSClient, err := tls.Get2WayTLSClient(config.CertsPath)
	if err != nil {
		defer log.Fatalf("fatal error setting up mutual TLS client")
	}
	shamirClient := shamir.NewSCClient(config.ShamirHost, mTLSClient)
	service := service.NewDistributionService(pmClient, eClient, r2pClient, shamirClient, db)

	switch distribute {
	case "advisories":
		log.Printf("Distributing to advisories")
		err := service.DistributeToAdvisoriesOnce()
		if err != nil {
			log.Printf("Error occurred during advisory distribution: %v", err)
		}
		return
	case "validators":
		log.Printf("Distributing to validators")
		service.DistributeToValidators()
		return
	default:
		// No distribution option specified, handle accordingly
		fmt.Println("No distribution option specified. Proceeding with default behavior.")
		if err = service.Run(); err != nil {
			log.Panicf("error occurred while spinning up service: %v", err)
		}
	}
}

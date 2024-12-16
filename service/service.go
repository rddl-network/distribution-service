package service

import (
	"context"
	"sync"

	"github.com/planetmint/planetmint-go/util"
	"github.com/rddl-network/distribution-service/config"
	log "github.com/rddl-network/go-logger"
	r2p "github.com/rddl-network/rddl-2-plmnt-service/client"
	shamir "github.com/rddl-network/shamir-coordinator-service/client"
	"github.com/robfig/cron/v3"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	wg sync.WaitGroup
)

type DistributionService struct {
	pmClient     IPlanetmintClient
	eClient      IElementsClient
	r2pClient    r2p.IR2PClient
	shamirClient shamir.ISCClient
	db           *leveldb.DB
	dbMutex      sync.Mutex
	logger       log.AppLogger
}

func NewDistributionService(pmClient IPlanetmintClient, eClient IElementsClient, r2pClient r2p.IR2PClient, shamirClient shamir.ISCClient, db *leveldb.DB) *DistributionService {
	cfg := config.GetConfig()
	service := &DistributionService{
		pmClient:     pmClient,
		eClient:      eClient,
		r2pClient:    r2pClient,
		shamirClient: shamirClient,
		db:           db,
		logger:       log.GetLogger(cfg.LogLevel),
	}
	_, err := service.ReadLastBlockHeight()
	if err != nil {
		err = service.WriteLastBlockHeight(0)
		if err != nil {
			service.logger.Error("error", "Cannot write Block height into file.")
		}
	}
	return service
}

// Run starts cronjob like thread to periodically check for DAO rewards to distribute to validators
func (ds *DistributionService) Run() (err error) {
	cfg := config.GetConfig()
	wg.Add(1)

	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)

	c := cron.New(cron.WithParser(parser), cron.WithChain())
	_, err = c.AddFunc(cfg.Cron, ds.DistributeToValidators)
	if err != nil {
		return
	}
	_, err = c.AddFunc(cfg.AdvisoryCron, ds.DistributeToAdvisories)
	if err != nil {
		return
	}
	c.Start()

	defer wg.Done()
	wg.Wait()

	return
}

func (ds *DistributionService) sendToAddresses(addresses []string, amount uint64, asset string) (err error) {
	amtString := util.UintValueToRDDLTokenString(amount)

	for _, address := range addresses {
		err = ds.sendAmountStringToAddress(address, amtString, asset)
		if err != nil {
			return
		}
	}

	return
}

func (ds *DistributionService) sendToAddress(address string, amount float64, asset string) (err error) {
	amountUint := util.RDDLToken2Uint(amount)
	amtString := util.UintValueToRDDLTokenString(amountUint)
	err = ds.sendAmountStringToAddress(address, amtString, asset)
	if err != nil {
		return
	}

	return
}

func (ds *DistributionService) sendAmountStringToAddress(address string, amount string, asset string) (err error) {
	_, err = ds.shamirClient.SendTokens(context.Background(), address, amount, asset)
	if err != nil {
		return
	}

	return
}

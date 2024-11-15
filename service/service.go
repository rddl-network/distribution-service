package service

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

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
	shamirClient shamir.IShamirCoordinatorClient
	db           *leveldb.DB
	dbMutex      sync.Mutex
	logger       log.AppLogger
}

func NewDistributionService(pmClient IPlanetmintClient, eClient IElementsClient, r2pClient r2p.IR2PClient, shamirClient shamir.IShamirCoordinatorClient, db *leveldb.DB) *DistributionService {
	cfg := config.GetConfig()

	return &DistributionService{
		pmClient:     pmClient,
		eClient:      eClient,
		r2pClient:    r2pClient,
		shamirClient: shamirClient,
		db:           db,
		logger:       log.GetLogger(cfg.LogLevel),
	}
}

// Run starts cronjob like thread to periodically check for DAO rewards to distribute to validators
func (ds *DistributionService) Run(cronExp string) (err error) {
	wg.Add(1)

	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)

	c := cron.New(cron.WithParser(parser), cron.WithChain())
	_, err = c.AddFunc(cronExp, ds.Distribute)
	if err != nil {
		return
	}
	c.Start()

	defer wg.Done()
	wg.Wait()

	return
}

// Distributes 10% of received funds to all validators
func (ds *DistributionService) Distribute() {
	distributionAmt, err := ds.getDistributionAmount()
	if err != nil {
		ds.logger.Error("msg", "Error while calculating distribution amount: "+err.Error())
		return
	}

	if distributionAmt == 0 {
		ds.logger.Error("msg", "No tokens to distribute.")
		return
	}

	liquidAddresses, err := ds.getBeneficiaries()
	if err != nil {
		ds.logger.Error("msg", "Error while fetching beneficiary addresses: "+err.Error())
		return
	}

	// CalculateShares
	share, _ := ds.calculateShares(distributionAmt, uint64(len(liquidAddresses)))

	// SendToAddresses
	ds.logger.Info("msg", "sending tokens", "addresses", strings.Join(liquidAddresses, ","), "amount", distributionAmt, "share", share)
	err = ds.sendToAddresses(share, liquidAddresses)
	if err != nil {
		ds.logger.Error("msg", "Error while sending to validators: "+err.Error())
		return
	}
}

func (ds *DistributionService) getDistributionAmount() (distributionAmt uint64, err error) {
	received, err := ds.checkReceivedBalance()
	if err != nil {
		return
	}

	ds.logger.Debug("msg", "Reading last occurrence")
	occurrence, err := ds.GetLastOccurrence()
	if err != nil {
		return
	}

	ds.logger.Debug("msg", "Storing current occurrence")
	err = ds.StoreOccurrence(time.Now().Unix(), received)
	if err != nil {
		return
	}

	if occurrence == nil {
		return CalculateDistributionAmount(0, received), nil
	}

	return CalculateDistributionAmount(occurrence.Amount, received), nil
}

// Checks for received asset on a given address
func (ds *DistributionService) checkReceivedBalance() (received uint64, err error) {
	cfg := config.GetConfig()
	ds.logger.Info("msg", "checking received balance", "address", cfg.FundAddress, " asset", cfg.Asset)

	confirmationString := strconv.Itoa(cfg.Confirmations)
	txDetails, err := ds.eClient.ListReceivedByAddress(cfg.GetElementsURL(),
		[]string{confirmationString, "false", "true", `"` + cfg.FundAddress + `"`, `"` + cfg.Asset + `"`},
	)
	if err != nil {
		return 0, err
	}

	for _, txDetail := range txDetails {
		received += util.RDDLToken2Uint(txDetail.Amount)
	}

	return
}

func (ds *DistributionService) getBeneficiaries() (addresses []string, err error) {
	plmntAddresses, err := ds.getActiveValidatorAddresses()
	if err != nil {
		return nil, err
	}

	ds.logger.Info("msg", "fetching liquid receive addresses", "planetmintAddresses", strings.Join(plmntAddresses, ","))
	return ds.getReceiveAddresses(plmntAddresses)
}

// getReceiveAddresses fetches receive addresses from the rddl-2-plmnt service
func (ds *DistributionService) getReceiveAddresses(addresses []string) (receiveAddresses []string, err error) {
	for _, address := range addresses {
		receiveAddress, err := ds.r2pClient.GetReceiveAddress(context.Background(), address)
		if err != nil {
			return nil, err
		}
		receiveAddresses = append(receiveAddresses, receiveAddress.LiquidAddress)
	}
	return
}

// Gets all active validator addresses
func (ds *DistributionService) getActiveValidatorAddresses() (addresses []string, err error) {
	valAddresses, err := ds.pmClient.GetValidatorAddresses()
	if err != nil {
		return nil, err
	}

	for _, address := range valAddresses {
		delegationAddresses, err := ds.pmClient.GetValidatorDelegationAddresses(address)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, delegationAddresses...)
	}

	return
}

// Calculates share per given address
func (ds *DistributionService) calculateShares(total uint64, numValidators uint64) (share uint64, remainder uint64) {
	if numValidators == 0 {
		return 0, total
	}

	share = total / numValidators
	remainder = total % numValidators
	return
}

func (ds *DistributionService) sendToAddresses(amount uint64, addresses []string) (err error) {
	amtString := util.UintValueToRDDLTokenString(amount)

	for _, address := range addresses {
		_, err = ds.shamirClient.SendTokens(context.Background(), address, amtString)
		if err != nil {
			return err
		}
	}

	return
}

func CalculateDistributionAmount(prev uint64, curr uint64) (distributionAmt uint64) {
	if prev == 0 {
		return curr / 100 * 10
	}

	return (curr - prev) / 100 * 10
}

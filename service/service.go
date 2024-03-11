package service

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/rddl-network/distribution-service/config"
	"github.com/robfig/cron/v3"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	wg sync.WaitGroup
)

type DistributionService struct {
	pmClient  IPlanetmintClient
	eClient   IElementsClient
	r2pClient IR2PClient
	db        *leveldb.DB
	dbMutex   sync.Mutex
}

func NewDistributionService(pmClient IPlanetmintClient, eClient IElementsClient, r2pClient IR2PClient, db *leveldb.DB) *DistributionService {
	return &DistributionService{
		pmClient:  pmClient,
		eClient:   eClient,
		r2pClient: r2pClient,
		db:        db,
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
	timestamp, err := ds.getLastOccurence()
	if err != nil {
		log.Println("Error while reading last occurence: " + err.Error())
		return
	}

	received := ds.CheckReceivedBalance(timestamp)

	// GetActiveValidatorAddresses
	plmntAddresses, err := ds.GetActiveValidatorAddresses()
	if err != nil {
		log.Println("Error while fetching validator set: " + err.Error())
		return
	}

	// CalculateShares
	share, _ := ds.CalculateShares(received, uint64(len(plmntAddresses)))

	liquidAddresses, err := ds.GetReceiveAddresses(plmntAddresses)
	if err != nil {
		log.Println("Error while fetching receive addresses: " + err.Error())
		return
	}

	// SendToAddresses
	err = ds.sendToAddresses(share, liquidAddresses)
	if err != nil {
		log.Println("Error while sending to validators: " + err.Error())
		return
	}

	err = ds.storeLastOccurence(time.Now().Unix())
	if err != nil {
		log.Println("Error while storing last occurence: " + err.Error())
	}
}

// Checks for Received RDDL since a given timestamp
func (ds *DistributionService) CheckReceivedBalance(timestamp int64) (received uint64) {
	return
}

// GetReceiveAddresses fetches receive addresses from the rddl-2-plmnt service
func (ds *DistributionService) GetReceiveAddresses(addresses []string) (receiveAddresses []string, err error) {
	for _, address := range addresses {
		receiveAddress, err := ds.r2pClient.GetReceiveAddress(address)
		if err != nil {
			return nil, err
		}
		receiveAddresses = append(receiveAddresses, receiveAddress)
	}
	return
}

// Gets all active validator addresses
func (ds *DistributionService) GetActiveValidatorAddresses() (addresses []string, err error) {
	valAddresses, err := ds.pmClient.GetValidatorAddresses()
	if err != nil {
		panic(err)
	}

	for _, address := range valAddresses {
		delegationAddresses, err := ds.pmClient.GetValidatorDelegationAddresses(address)
		if err != nil {
			panic(err)
		}
		addresses = append(addresses, delegationAddresses...)
	}

	return
}

// Calculates share per given address
func (ds *DistributionService) CalculateShares(total uint64, numValidators uint64) (share uint64, remainder uint64) {
	if numValidators == 0 {
		return 0, total
	}

	share = total / numValidators
	remainder = total % numValidators
	return
}

func (ds *DistributionService) sendToAddresses(amount uint64, addresses []string) (err error) {
	for _, address := range addresses {
		err = ds.issueShamirTransaction(amount, address)
		if err != nil {
			return err
		}
	}

	return
}

func (ds *DistributionService) issueShamirTransaction(amount uint64, address string) (err error) {
	cfg := config.GetConfig()
	url := fmt.Sprintf("%s/%s/%d", cfg.ShamireHost, address, amount)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return
}

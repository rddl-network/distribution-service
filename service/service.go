package service

import (
	"log"
	"sync"
	"time"

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
	addresses, err := ds.GetActiveValidatorAddresses()
	if err != nil {
		log.Println("Error while fetching validator set: " + err.Error())
		return
	}

	// CalculateShares
	share, _ := ds.CalculateShares(received, uint64(len(addresses)))

	// SendToAddresses
	err = ds.SendToAddresses(share, addresses)
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

func (ds *DistributionService) SendToAddresses(share uint64, addresses []string) (err error) {
	return
}

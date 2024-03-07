package service

import (
	"fmt"
	"sync"

	"github.com/robfig/cron/v3"
)

var (
	wg sync.WaitGroup
)

type DistributionService struct {
	pmClient  IPlanetmintClient
	eClient   IElementsClient
	r2pClient IR2PClient
}

func NewDistributionService(pmClient IPlanetmintClient, eClient IElementsClient, r2pClient IR2PClient) *DistributionService {
	return &DistributionService{
		pmClient:  pmClient,
		eClient:   eClient,
		r2pClient: r2pClient,
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
	fmt.Println("CRON THE SECRET OF THE TASK")
}

// Checks for Received RDDL over a given timeperiod
func (ds *DistributionService) CheckReceivedBalance() {}

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

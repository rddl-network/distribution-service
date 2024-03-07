package service

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

// Start service/ticker to periodically check for DAO rewards to distribute to validators
// TODO: Make configurable by cron job like syntax
func (ds *DistributionService) Run() (err error) {
	return
}

// Distributes 10% of received funds to all validators
func (ds *DistributionService) Distribute() {}

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

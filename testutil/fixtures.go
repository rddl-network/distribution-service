package testutil

import elementstypes "github.com/rddl-network/elements-rpc/types"

var TxDetails = []elementstypes.ListReceivedByAddressResult{
	{
		Amount: 100,
	},
}

var TxZeroDetails = []elementstypes.ListReceivedByAddressResult{
	{
		Amount: 0,
	},
}

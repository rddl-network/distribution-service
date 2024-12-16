package service_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rddl-network/distribution-service/service"
	"github.com/rddl-network/distribution-service/testutil"
	"github.com/stretchr/testify/assert"

	r2p "github.com/rddl-network/rddl-2-plmnt-service/types"
	shamir "github.com/rddl-network/shamir-coordinator-service/types"
)

func TestValidatorDistribution(t *testing.T) {
	app, db, mocks := setupService(t)
	defer db.Close()

	mocks.eClientMock.EXPECT().ListReceivedByAddress(gomock.Any(), gomock.Any()).Times(1).Return(testutil.TxDetails, nil)
	mocks.pmClientMock.EXPECT().GetValidatorAddresses().Times(1).Return([]string{"valoper1", "valoper2"}, nil)
	mocks.pmClientMock.EXPECT().GetValidatorDelegationAddresses("valoper1").Times(1).Return([]string{"val1"}, nil)
	mocks.pmClientMock.EXPECT().GetValidatorDelegationAddresses("valoper2").Times(1).Return([]string{"val2"}, nil)
	mocks.r2pClientMock.EXPECT().GetReceiveAddress(gomock.Any(), "val1").Times(1).Return(r2p.ReceiveAddressResponse{LiquidAddress: "liquid1"}, nil)
	mocks.r2pClientMock.EXPECT().GetReceiveAddress(gomock.Any(), "val2").Times(1).Return(r2p.ReceiveAddressResponse{LiquidAddress: "liquid2"}, nil)
	mocks.shamirClientMock.EXPECT().SendTokens(gomock.Any(), "liquid1", "5.00000000", gomock.Any()).Times(1).Return(shamir.SendTokensResponse{}, nil)
	mocks.shamirClientMock.EXPECT().SendTokens(gomock.Any(), "liquid2", "5.00000000", gomock.Any()).Times(1).Return(shamir.SendTokensResponse{}, nil)

	app.DistributeToValidators()
}

func TestServiceZeroVaildatorDistribution(t *testing.T) {
	app, db, mocks := setupService(t)
	defer db.Close()

	mocks.eClientMock.EXPECT().ListReceivedByAddress(gomock.Any(), gomock.Any()).Times(1).Return(testutil.TxZeroDetails, nil)
	// the other mocks are not called : sanity check
	mocks.pmClientMock.EXPECT().GetValidatorAddresses().Times(0)
	mocks.pmClientMock.EXPECT().GetValidatorDelegationAddresses("valoper1").Times(0)
	mocks.pmClientMock.EXPECT().GetValidatorDelegationAddresses("valoper2").Times(0)
	mocks.r2pClientMock.EXPECT().GetReceiveAddress(gomock.Any(), "val1").Times(0)
	mocks.r2pClientMock.EXPECT().GetReceiveAddress(gomock.Any(), "val2").Times(0)
	mocks.shamirClientMock.EXPECT().SendTokens(gomock.Any(), "liquid1", "5.00000000", gomock.Any()).Times(0)
	mocks.shamirClientMock.EXPECT().SendTokens(gomock.Any(), "liquid2", "5.00000000", gomock.Any()).Times(0)
	app.DistributeToValidators()
}

// Using uint64 with at least 8 zeros appended indicating the shift from float string to uint representation
func TestCalculateVaildatorDistributionAmount(t *testing.T) {
	amt := service.CalculateValidatorDistributionAmount(100000000, 1000000000)
	assert.Equal(t, uint64(90000000), amt)

	amt = service.CalculateValidatorDistributionAmount(0, 1000000000)
	assert.Equal(t, uint64(100000000), amt)
}

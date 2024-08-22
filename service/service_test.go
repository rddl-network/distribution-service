package service_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rddl-network/distribution-service/service"
	"github.com/rddl-network/distribution-service/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/storage"

	r2p "github.com/rddl-network/rddl-2-plmnt-service/service"
	shamir "github.com/rddl-network/shamir-coordinator-service/service"
)

type Mocks struct {
	eClientMock      *testutil.MockIElementsClient
	pmClientMock     *testutil.MockIPlanetmintClient
	r2pClientMock    *testutil.MockIR2PClient
	shamirClientMock *testutil.MockIShamirCoordinatorClient
}

func setupService(t *testing.T) (app *service.DistributionService, db *leveldb.DB, mocks Mocks) {
	db, err := leveldb.Open(storage.NewMemStorage(), nil)
	if err != nil {
		t.Fatal("Error opening in-memory LevelDB: ", err.Error())
	}

	ctrl := gomock.NewController(t)
	eClientMock := testutil.NewMockIElementsClient(ctrl)
	pmClientMock := testutil.NewMockIPlanetmintClient(ctrl)
	r2pClientMock := testutil.NewMockIR2PClient(ctrl)
	shamirClientMock := testutil.NewMockIShamirCoordinatorClient(ctrl)

	app = service.NewDistributionService(pmClientMock, eClientMock, r2pClientMock, shamirClientMock, db)

	mocks.eClientMock = eClientMock
	mocks.pmClientMock = pmClientMock
	mocks.r2pClientMock = r2pClientMock
	mocks.shamirClientMock = shamirClientMock

	return
}

func TestService(t *testing.T) {
	app, db, mocks := setupService(t)
	defer db.Close()

	mocks.eClientMock.EXPECT().ListReceivedByAddress(gomock.Any(), gomock.Any()).Times(1).Return(testutil.TxDetails, nil)
	mocks.pmClientMock.EXPECT().GetValidatorAddresses().Times(1).Return([]string{"valoper1", "valoper2"}, nil)
	mocks.pmClientMock.EXPECT().GetValidatorDelegationAddresses("valoper1").Times(1).Return([]string{"val1"}, nil)
	mocks.pmClientMock.EXPECT().GetValidatorDelegationAddresses("valoper2").Times(1).Return([]string{"val2"}, nil)
	mocks.r2pClientMock.EXPECT().GetReceiveAddress(gomock.Any(), "val1").Times(1).Return(r2p.ReceiveAddressResponse{LiquidAddress: "liquid1"}, nil)
	mocks.r2pClientMock.EXPECT().GetReceiveAddress(gomock.Any(), "val2").Times(1).Return(r2p.ReceiveAddressResponse{LiquidAddress: "liquid2"}, nil)
	mocks.shamirClientMock.EXPECT().SendTokens(gomock.Any(), "liquid1", "5.00000000").Times(1).Return(shamir.SendTokensResponse{}, nil)
	mocks.shamirClientMock.EXPECT().SendTokens(gomock.Any(), "liquid2", "5.00000000").Times(1).Return(shamir.SendTokensResponse{}, nil)

	app.Distribute()
}

func TestServiceZeroDistribution(t *testing.T) {
	app, db, mocks := setupService(t)
	defer db.Close()

	mocks.eClientMock.EXPECT().ListReceivedByAddress(gomock.Any(), gomock.Any()).Times(1).Return(testutil.TxZeroDetails, nil)
	// the other mocks are not called : sanity check
	mocks.pmClientMock.EXPECT().GetValidatorAddresses().Times(0)
	mocks.pmClientMock.EXPECT().GetValidatorDelegationAddresses("valoper1").Times(0)
	mocks.pmClientMock.EXPECT().GetValidatorDelegationAddresses("valoper2").Times(0)
	mocks.r2pClientMock.EXPECT().GetReceiveAddress(gomock.Any(), "val1").Times(0)
	mocks.r2pClientMock.EXPECT().GetReceiveAddress(gomock.Any(), "val2").Times(0)
	mocks.shamirClientMock.EXPECT().SendTokens(gomock.Any(), "liquid1", "5.00000000").Times(0)
	mocks.shamirClientMock.EXPECT().SendTokens(gomock.Any(), "liquid2", "5.00000000").Times(0)
	app.Distribute()
}

// Using uint64 with at least 8 zeros appended indicating the shift from float string to uint representation
func TestCalculateDistributionAmount(t *testing.T) {
	amt := service.CalculateDistributionAmount(100000000, 1000000000)
	assert.Equal(t, uint64(90000000), amt)

	amt = service.CalculateDistributionAmount(0, 1000000000)
	assert.Equal(t, uint64(100000000), amt)
}

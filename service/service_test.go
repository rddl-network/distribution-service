package service_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rddl-network/distribution-service/service"
	"github.com/rddl-network/distribution-service/testutil"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

type Mocks struct {
	eClientMock      *testutil.MockIElementsClient
	pmClientMock     *testutil.MockIPlanetmintClient
	r2pClientMock    *testutil.MockIR2PClient
	shamirClientMock *testutil.MockIShamirClient
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
	shamirClientMock := testutil.NewMockIShamirClient(ctrl)

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
	mocks.r2pClientMock.EXPECT().GetReceiveAddress("val1").Times(1).Return("liquid1", nil)
	mocks.r2pClientMock.EXPECT().GetReceiveAddress("val2").Times(1).Return("liquid2", nil)
	mocks.shamirClientMock.EXPECT().IssueTransaction("5.00000000", "liquid1").Times(1).Return(nil)
	mocks.shamirClientMock.EXPECT().IssueTransaction("5.00000000", "liquid2").Times(1).Return(nil)

	app.Distribute()
}

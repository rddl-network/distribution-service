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
	shamirClientMock *testutil.MockISCClient
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
	shamirClientMock := testutil.NewMockISCClient(ctrl)

	app = service.NewDistributionService(pmClientMock, eClientMock, r2pClientMock, shamirClientMock, db)

	mocks.eClientMock = eClientMock
	mocks.pmClientMock = pmClientMock
	mocks.r2pClientMock = r2pClientMock
	mocks.shamirClientMock = shamirClientMock

	return
}

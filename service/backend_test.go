package service_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rddl-network/distribution-service/service"
	"github.com/rddl-network/distribution-service/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

func setupService(t *testing.T) (app *service.DistributionService, db *leveldb.DB) {
	db, err := leveldb.Open(storage.NewMemStorage(), nil)
	if err != nil {
		t.Fatal("Error opening in-memory LevelDB: ", err.Error())
	}

	ctrl := gomock.NewController(t)
	eClient := testutil.NewMockIElementsClient(ctrl)
	pmClient := testutil.NewMockIPlanetmintClient(ctrl)
	r2pClient := testutil.NewMockIR2PClient(ctrl)

	app = service.NewDistributionService(pmClient, eClient, r2pClient, db)
	return
}

func createNOccurences(app *service.DistributionService, n int) []service.Occurence {
	items := make([]service.Occurence, n)
	for i := range items {
		items[i].Timestamp = int64(i)
		items[i].Amount = uint64(i * 1000)
		app.StoreOccurence(items[i].Timestamp, items[i].Amount)
	}
	return items
}

func TestOccurence(t *testing.T) {
	app, db := setupService(t)
	defer db.Close()

	items := createNOccurences(app, 10)
	occurence, err := app.GetOccurence(items[5].Timestamp)
	assert.NoError(t, err)
	assert.Equal(t, items[5].Amount, occurence.Amount)

	lastOccurence, err := app.GetLastOccurence()
	assert.NoError(t, err)
	assert.Equal(t, items[9].Amount, lastOccurence.Amount)
}

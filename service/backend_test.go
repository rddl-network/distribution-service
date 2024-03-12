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

func createNOccurrences(app *service.DistributionService, n int) []service.Occurrence {
	items := make([]service.Occurrence, n)
	for i := range items {
		items[i].Timestamp = int64(i)
		items[i].Amount = uint64(i * 1000)
		_ = app.StoreOccurrence(items[i].Timestamp, items[i].Amount)
	}
	return items
}

func TestOccurrence(t *testing.T) {
	app, db := setupService(t)
	defer db.Close()

	items := createNOccurrences(app, 10)
	occurrence, err := app.GetOccurrence(items[5].Timestamp)
	assert.NoError(t, err)
	assert.Equal(t, items[5].Amount, occurrence.Amount)

	lastOccurrence, err := app.GetLastOccurrence()
	assert.NoError(t, err)
	assert.Equal(t, items[9].Amount, lastOccurrence.Amount)
}

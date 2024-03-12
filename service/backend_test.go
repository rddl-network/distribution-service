package service_test

import (
	"testing"

	"github.com/rddl-network/distribution-service/service"
	"github.com/stretchr/testify/assert"
)

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
	app, db, _ := setupService(t)
	defer db.Close()

	items := createNOccurrences(app, 10)
	occurrence, err := app.GetOccurrence(items[5].Timestamp)
	assert.NoError(t, err)
	assert.Equal(t, items[5].Amount, occurrence.Amount)

	lastOccurrence, err := app.GetLastOccurrence()
	assert.NoError(t, err)
	assert.Equal(t, items[9].Amount, lastOccurrence.Amount)
}

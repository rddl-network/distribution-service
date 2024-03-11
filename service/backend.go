package service

import (
	"encoding/json"
)

var (
	key = []byte("lastoccurence")
)

type Occurence struct {
	Timestamp   int64  `json:"timestamp"`
	TotalAmount uint64 `json:"total-amount"`
}

func (ds *DistributionService) storeLastOccurence(timestamp int64, amount uint64) (err error) {
	occurence, err := ds.getLastOccurence()
	if err != nil {
		return err
	}

	occurence = &Occurence{
		Timestamp:   timestamp,
		TotalAmount: occurence.TotalAmount + amount,
	}

	bytes, err := json.Marshal(occurence)
	if err != nil {
		return err
	}

	ds.dbMutex.Lock()
	err = ds.db.Put(key, bytes, nil)
	ds.dbMutex.Unlock()

	return
}

func (ds *DistributionService) getLastOccurence() (occurence *Occurence, err error) {
	bytes, err := ds.db.Get(key, nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &occurence)
	if err != nil {
		return nil, err
	}

	return occurence, nil
}

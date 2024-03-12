package service

import (
	"encoding/binary"
	"encoding/json"

	"github.com/syndtr/goleveldb/leveldb/util"
)

const (
	OccurencePrefix = "Occurence/"
)

func occurenceKey(timestamp int64) (key []byte) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(timestamp))

	prefixBytes := []byte(OccurencePrefix)
	key = append(key, prefixBytes...)
	key = append(key, buf...)

	return
}

type Occurence struct {
	Timestamp int64  `json:"timestamp"`
	Amount    uint64 `json:"amount"`
}

func (ds *DistributionService) StoreOccurence(timestamp int64, amount uint64) (err error) {
	occurence := &Occurence{
		Timestamp: timestamp,
		Amount:    amount,
	}

	bytes, err := json.Marshal(occurence)
	if err != nil {
		return err
	}

	ds.dbMutex.Lock()
	err = ds.db.Put(occurenceKey(timestamp), bytes, nil)
	ds.dbMutex.Unlock()

	return
}

func (ds *DistributionService) GetOccurence(timestamp int64) (occurence *Occurence, err error) {
	bytes, err := ds.db.Get(occurenceKey(timestamp), nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &occurence)
	if err != nil {
		return nil, err
	}

	return occurence, nil
}

func (ds *DistributionService) GetLastOccurence() (occurence *Occurence, err error) {
	iter := ds.db.NewIterator(util.BytesPrefix([]byte(OccurencePrefix)), nil)
	defer iter.Release()

	iter.Last()

	bytes := iter.Value()
	err = json.Unmarshal(bytes, &occurence)
	return
}

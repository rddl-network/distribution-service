package service

import (
	"encoding/binary"
	"encoding/json"

	"github.com/syndtr/goleveldb/leveldb/util"
)

const (
	OccurrencePrefix = "Occurrence/"
)

func occurrenceKey(timestamp int64) (key []byte) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(timestamp))

	prefixBytes := []byte(OccurrencePrefix)
	key = append(key, prefixBytes...)
	key = append(key, buf...)

	return
}

type Occurrence struct {
	Timestamp int64  `json:"timestamp"`
	Amount    uint64 `json:"amount"`
}

func (ds *DistributionService) StoreOccurrence(timestamp int64, amount uint64) (err error) {
	occurrence := &Occurrence{
		Timestamp: timestamp,
		Amount:    amount,
	}

	bytes, err := json.Marshal(occurrence)
	if err != nil {
		return err
	}

	ds.dbMutex.Lock()
	err = ds.db.Put(occurrenceKey(timestamp), bytes, nil)
	ds.dbMutex.Unlock()

	return
}

func (ds *DistributionService) GetOccurrence(timestamp int64) (occurrence *Occurrence, err error) {
	bytes, err := ds.db.Get(occurrenceKey(timestamp), nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &occurrence)
	if err != nil {
		return nil, err
	}

	return occurrence, nil
}

func (ds *DistributionService) GetLastOccurrence() (occurrence *Occurrence, err error) {
	iter := ds.db.NewIterator(util.BytesPrefix([]byte(OccurrencePrefix)), nil)
	defer iter.Release()

	iter.Last()

	bytes := iter.Value()
	err = json.Unmarshal(bytes, &occurrence)
	return
}

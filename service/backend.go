package service

import "encoding/binary"

var (
	key = []byte("lastoccurence")
)

func (ds *DistributionService) storeLastOccurence(timestamp int64) (err error) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(timestamp))

	ds.dbMutex.Lock()
	err = ds.db.Put(key, b, nil)
	ds.dbMutex.Unlock()

	return
}

func (ds *DistributionService) getLastOccurence() (timestamp int64, err error) {
	bytes, err := ds.db.Get(key, nil)

	timestamp = int64(binary.BigEndian.Uint64(bytes))

	return
}

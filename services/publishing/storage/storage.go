package storage

import (
	mdb "github.com/shawntoffel/GoMongoDb"
)

type Storage struct {
	MessageStorage        MessageStorage
	MessageRequestStorage MessageRequestStorage
}

type StorageConfig struct {
	MessageStorageConfig        mdb.DbConfig
	MessageRequestStorageConfig mdb.DbConfig
}

func NewStorage(storageConfig StorageConfig) (Storage, error) {

	storage := Storage{}

	var err error

	storage.MessageStorage, err = NewMessageStorage(storageConfig.MessageStorageConfig)

	if err != nil {
		return storage, err
	}

	storage.MessageRequestStorage, err = NewMessageRequestStorage(storageConfig.MessageRequestStorageConfig)

	return storage, err
}

func (s *Storage) Close() {
	s.MessageStorage.Close()
	s.MessageRequestStorage.Close()
}

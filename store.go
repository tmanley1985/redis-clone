package main

import "errors"

// TODO: I need to separate these out into their own packages.

type DataStoreInterface interface {
	Get(key string) (string, bool)
	Set(key, value string) error
}

type ClientResponse struct {
	Value interface{}
	Err   error
}

type GetOperation struct {
	Key      string
	Response chan ClientResponse
}

type SetOperation struct {
	Key      string
	Value    string
	Response chan ClientResponse
}

type DataStore struct {
	getChan  chan GetOperation
	setChan  chan SetOperation
	data     map[string]string
	shutdown chan struct{}
}

type DataStoreConfig struct {
	ChannelBuffer   int
	PersistencePath string
}

func NewDataStore(config DataStoreConfig) (*DataStore, error) {
	ds := &DataStore{
		getChan:  make(chan GetOperation, config.ChannelBuffer),
		setChan:  make(chan SetOperation, config.ChannelBuffer),
		data:     make(map[string]string),
		shutdown: make(chan struct{}),
	}

	if config.PersistencePath != "" {
		err := ds.loadFromPersistence(config.PersistencePath)
		if err != nil {
			return nil, err
		}
	}

	// Maybe I should remove this and let the caller call this? I don't like the implicit behavior.
	go ds.run()
	return ds, nil
}

func (ds *DataStore) loadFromPersistence(path string) error {
	// Here I would need to read from an append only log that I have not yet implemented.
	return nil
}

func (ds *DataStore) run() {
	for {
		select {
		case getOp := <-ds.getChan:
			value, exists := ds.data[getOp.Key]
			if exists {
				getOp.Response <- ClientResponse{Value: value}
			} else {
				getOp.Response <- ClientResponse{Err: errors.New("key not found")}
			}
		case setOp := <-ds.setChan:
			ds.data[setOp.Key] = setOp.Value
			setOp.Response <- ClientResponse{Value: "OK"}
		case <-ds.shutdown:
			// Perform any cleanup here
			return
		}
	}
}

func (ds *DataStore) Shutdown() {
	close(ds.shutdown)
}

func (ds *DataStore) Get(key string, responseChan chan ClientResponse) {
	ds.getChan <- GetOperation{Key: key, Response: responseChan}
}

func (ds *DataStore) Set(key, value string, responseChan chan ClientResponse) {
	ds.setChan <- SetOperation{Key: key, Value: value, Response: responseChan}
}

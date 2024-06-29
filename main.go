package main

import (
	"log"
)

func main() {

	dataStore, err := NewDataStore(DataStoreConfig{})

	if err != nil {
		log.Fatal(err)
	}

	server := NewServer(Config{}, dataStore)
	log.Fatal(server.Start())
}

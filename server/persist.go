package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"time"
)

type SerializableStore struct {
	Data map[string]string
	TTL  map[string]time.Time
}

type Persister struct {
	Filename string
}

func NewPersister(filename string) *Persister {
	return &Persister{Filename: filename}
}

func (p *Persister) Save(store SerializableStore) error {
	file, err := os.Create(p.Filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(store)
}

func (p *Persister) Load() (SerializableStore, error) {
	file, err := os.Open(p.Filename)
	if err != nil {
		fmt.Println("HERE?")
		return SerializableStore{}, err
	}
	defer file.Close()

	var store SerializableStore
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&store)
	return store, err
}

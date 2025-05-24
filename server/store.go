package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Store struct {
	mu      sync.RWMutex
	data    map[string]string
	ttl     map[string]time.Time
	persist bool
}

type StoreConfig struct {
	Persist          bool
	SnapshotFile     string
	SnapshotInterval time.Duration
}

func NewStore(config StoreConfig) *Store {
	store := &Store{
		data:    make(map[string]string),
		ttl:     make(map[string]time.Time),
		persist: config.Persist,
	}
	fmt.Println("%s", config.SnapshotFile)
	persister := NewPersister(config.SnapshotFile)

	if config.Persist {
		if loaded, err := persister.Load(); err == nil {
			store.data = loaded.Data
			store.ttl = loaded.TTL
			fmt.Println("%s", loaded.Data)
		}

		go store.startSnapshotting(persister, config.SnapshotInterval)
	}

	return store
}

func (s *Store) startSnapshotting(p *Persister, interval time.Duration) {
	ticker := time.NewTicker(interval)

	for range ticker.C {
		s.mu.RLock()
		dataCopy := make(map[string]string, len(s.data))
		ttlCopy := make(map[string]time.Time, len(s.data))

		for k, v := range s.data {
			dataCopy[k] = v
			ttlCopy[k] = s.ttl[k]
		}
		s.mu.RUnlock()
		if err := p.Save(SerializableStore{
			Data: dataCopy,
			TTL:  ttlCopy,
		}); err != nil {
			log.Println("Error saving snapshot:", err)
		}
	}
}

func (s *Store) Set(key, value string, expiry *time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
	if expiry != nil {
		s.ttl[key] = *expiry
	} else {
		delete(s.ttl, key)
	}
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, ok := s.data[key]
	if !ok {
		return "", false
	}

	if expiry, hasExpiry := s.ttl[key]; hasExpiry {
		if time.Now().After(expiry) {
			delete(s.data, key)
			delete(s.ttl, key)
			return "", false
		}
	}

	return val, true
}

func (s *Store) Del(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, ok := s.data[key]
	delete(s.data, key)
	delete(s.ttl, key)
	return val, ok
}

func (s *Store) GetAll() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	var keys []string

	for k := range s.data {
		if expiry, hasExpiry := s.ttl[k]; hasExpiry {
			if time.Now().After(expiry) {
				delete(s.data, k)
				delete(s.ttl, k)
				continue
			}
		}
		keys = append(keys, k)
	}

	return keys
}

func (s *Store) DelAll() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for k := range s.data {
		delete(s.data, k)
		delete(s.ttl, k)
	}
}

package jackdb

import (
	"strconv"
	"strings"
	"sync"
)

// Store structure
type Store struct {
	dataMap map[string]string
	subMap  map[string][]chan<- string
	sync.RWMutex
}

// Init initializes dataMap (which keeps the key-value records) and subMap (which keeps the subscription list on each key).
func (s *Store) Init() {
	s.dataMap = make(map[string]string)
	s.subMap = make(map[string][]chan<- string)
}

// Get value
func (s *Store) Get(key string) (value string, ok bool) {
	s.RLock()
	defer s.RUnlock()
	value, ok = s.dataMap[key]
	return
}

// Set key to value
func (s *Store) Set(key, value string) bool {
	s.Lock()
	defer s.Unlock()
	s.dataMap[key] = value
	return true
}

// Delete a key
func (s *Store) Delete(keys []string) {
	s.Lock()
	defer s.Unlock()
	for _, key := range keys {
		delete(s.dataMap, key)
	}
}

// Publish a stream to a key
func (s *Store) Publish(key string, incoming <-chan string) {
	for {
		value, ok := <-incoming
		if !ok {
			return
		}
		_ = s.Set(key, value)
		s.updateSubscribers(key, value)
	}
}

// Subscribe to published changes on a key
func (s *Store) Subscribe(key string, outgoing chan<- string) {
	_, hasSubs := s.fetchSubscribers(key)
	s.Lock()
	defer s.Unlock()
	if hasSubs {
		s.subMap[key] = append(s.subMap[key], outgoing)
	} else {
		subs := []chan<- string{outgoing}
		s.subMap[key] = subs
	}
}

// Unsubscribe to published changes on a key
func (s *Store) unsubscribe(key string, outgoing chan<- string) {
	subs, hasSubs := s.fetchSubscribers(key)
	s.Lock()
	defer s.Unlock()
	if hasSubs {
		newSubs := []chan<- string{}
		for _, sub := range subs {
			if sub == outgoing {
				continue
			}
			newSubs = append(newSubs, sub)
		}
		s.subMap[key] = newSubs
	}
}

// Unexported
func (s *Store) fetchSubscribers(key string) ([]chan<- string, bool) {
	s.RLock()
	subs, hasSubs := s.subMap[key]
	s.RUnlock()
	return subs, hasSubs
}

// Unexported
func (s *Store) updateSubscribers(key, value string) {
	subs, ok := s.fetchSubscribers(key)
	if ok {
		for _, out := range subs {
			defer func(o chan<- string) {
				if r := recover(); r != nil {
					s.unsubscribe(key, o)
				}
			}(out)
			out <- value
		}
	}
}

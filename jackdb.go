package main

import (
	"net"
	"strconv"
	"strings"
	"sync"
	"log"
	"flag"
)

// Store structure
type Store struct {
	dataMap map[string]string
	subMap  map[string][]chan<- string
	sync.RWMutex
}

// Initialize storage
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

var (
	DB       Store
	OnlyOnce sync.Once
)

var port = flag.Int("port", 2000, "tcp port number")

// Unexported
func initialize() {
	DB.Init()
}

// Main method
func main() {
	flag.Parse()
	StartServer(*port)
}

// Start server on a port
func StartServer(pt int) {
	OnlyOnce.Do(initialize)

	portStr := ":" + strconv.Itoa(pt)
	listener, err := net.Listen("tcp", portStr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		c, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go HandleConnection(c)
	}
}

// Close connection
func CloseConnection(c net.Conn) {
	log.Printf("[%s] closed connection\n", c.RemoteAddr())
	c.Close()
}

// Handling connection
func HandleConnection(c net.Conn) {
	defer CloseConnection(c)

	fromAddr := c.RemoteAddr()
	log.Printf("[%s] new connection\n", fromAddr)

	buf := make([]byte, 1024)

NEXTMESSAGE:

	nb, err := c.Read(buf)
	if err != nil {
		return
	}

	args := strings.Split(string(buf[:nb]), " ")
	numArgs := len(args)

	switch args[0] {

	case "GET":

		if numArgs != 2 {
			return
		}

		if value, ok := DB.Get(args[1]); ok {
			_, err = c.Write([]byte(value))
		} else {
			_, err = c.Write([]byte("(nil)"))
		}

		if err != nil {
			return
		}

		goto NEXTMESSAGE

	case "SET":

		if numArgs != 3 {
			return
		}

		if ok := DB.Set(args[1], args[2]); ok {
			_, err = c.Write([]byte("OK"))
		} else {
			_, err = c.Write([]byte("FAIL"))
			return
		}

		if err != nil {
			return
		}

		goto NEXTMESSAGE

	case "DEL":

		if numArgs < 2 {
			return
		}

		DB.Delete(args[1:])
		_, err = c.Write([]byte("OK"))

		if err != nil {
			return
		}

		goto NEXTMESSAGE

	case "PUB":

		if numArgs != 2 {
			return
		}

		incoming := make(chan string)
		go DB.Publish(args[1], incoming)

		_, err = c.Write([]byte("READY"))

		if err != nil {
			close(incoming)
			return
		}

		for {
			nb, err := c.Read(buf)

			if err != nil {
				close(incoming)
				return
			}

			incoming <- string(buf[:nb])
		}

	case "SUB":

		if numArgs != 2 {
			return
		}

		outgoing := make(chan string)
		DB.Subscribe(args[1], outgoing)

		for value := range outgoing {
			_, err := c.Write([]byte(value))

			if err != nil {
				close(outgoing)
				return
			}
		}

	default:
		return
	}
}

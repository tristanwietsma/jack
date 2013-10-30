/*
Copyright 2013 Tristan Wietsma

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package jack

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
)

type connectionPoolError struct {
	desc string
}

func (e connectionPoolError) Error() string {
	return fmt.Sprintf("Connection Pool Error: %s", e.desc)
}

var maxConnectionsError = connectionPoolError{"Maximum connections reached."}

// The ConnectionPool struct maintains a finite number of connections for use on client-side.
type ConnectionPool struct {
	address string
	port    uint
	size    uint
	count   uint
	free    []*Connection
}

// NewConnectionPool constructs a ConnectionPool for a given address, port number, and pool size.
func NewConnectionPool(address string, port, size uint) *ConnectionPool {
	return &ConnectionPool{
		address: address,
		port:    port,
		size:    size,
	}
}

// Connect gets a connection from the pool
func (cp *ConnectionPool) Connect() (*Connection, error) {

	if cp.count == cp.size && len(cp.free) == 0 {
		return nil, maxConnectionsError
	}

	if len(cp.free) > 0 {
		sc := cp.free[0]
		cp.free = cp.free[1:]
		return sc, nil
	}

	cp.count++
	sc, err := NewConnection(cp.address, cp.port)
	if err != nil {
		return nil, err
	}
	return sc, nil
}

// Free sends a connection back to the pool
func (cp *ConnectionPool) Free(c *Connection) error {
	cp.free = append(cp.free, c)
	return nil
}

// The Connection struct defines an individual client connection.
type Connection struct {
	conn net.Conn
	feed chan string
}

func (sc *Connection) transmit(m *Message) {
	b := m.Bytes()
	_, err := sc.conn.Write(b)
	buf := make([]byte, 1024)

WAIT_FOR_SERVER:

	_, err = sc.conn.Read(buf)
	if err != nil {
		panic(err)
	}

	end := bytes.IndexByte(buf, EOM)
	if end < 0 {
		panic(ProtocolError{"Message is missing EOM byte."})
	}

	payload := string(buf[:end])
	sc.feed <- payload

	if m.cmd == SUB {
		goto WAIT_FOR_SERVER
	}
}

// Get returns the value associated with a given key.
func (sc *Connection) Get(key string) string {
	m := NewGetMessage(key)
	go sc.transmit(m)
	return <-sc.feed
}

// Set assigns a value to a key.
func (sc *Connection) Set(key, value string) string {
	m := NewSetMessage(key, value)
	go sc.transmit(m)
	return <-sc.feed
}

// Delete removes a key-value pair from the database.
func (sc *Connection) Delete(key string) string {
	m := NewDeleteMessage(key)
	go sc.transmit(m)
	return <-sc.feed
}

// Publish sets a value to a key and triggers a subscriber update.
func (sc *Connection) Publish(key, value string) string {
	m := NewPublishMessage(key, value)
	go sc.transmit(m)
	return <-sc.feed
}

// Subscribe added a channel to the subscription list on a key.
func (sc *Connection) Subscribe(key string, recv chan<- string) {
	m := NewSubscribeMessage(key)
	go sc.transmit(m)
	for {
		recv <- <-sc.feed
	}
}

// Close closes a connection.
func (sc *Connection) Close() error {
	err := sc.conn.Close()
	return err
}

// NewConnection returns a connection to a given address and port number.
func NewConnection(address string, port uint) (*Connection, error) {
	fullAddress := address + ":" + strconv.FormatUint(uint64(port), 10)
	conn, err := net.Dial("tcp", fullAddress)
	sc := Connection{}
	if err == nil {
		sc.conn = conn
	}
	sc.feed = make(chan string)
	return &sc, err
}

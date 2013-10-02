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

package jackdb

import (
	"fmt"
	"bytes"
	"net"
	"strconv"
)

type ConnectionPoolError struct {
	desc string
}

func (e ConnectionPoolError) Error() string {
	return fmt.Sprintf("ConnectionPoolError: %s", e.desc)
}

var MaxConnectionsError = ConnectionPoolError{"Maximum connections reached."}

// ConnectionPool
type ConnectionPool struct {
	address string
	port uint
	size uint
	count uint
	free []*Connection
}

func NewConnectionPool(address string, port, size uint) *ConnectionPool {
	return &ConnectionPool{
		address: address,
		port: port,
		size: size,
	}
}

// Connect gets a connection from the pool
func (cp *ConnectionPool) Connect() (*Connection, error) {

	if cp.count == cp.size && len(cp.free) == 0 {
		return nil, MaxConnectionsError
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

// Connection
type Connection struct {
	conn net.Conn
	feed chan string
}

func (sc *Connection) transmit(m *Message) error {
	b := m.Bytes()
	_, err := sc.conn.Write(b)
	buf := make([]byte, 1024)

WAIT_FOR_SERVER:

	fmt.Println("debug::: 93")

	_, err = sc.conn.Read(buf)
	if err != nil {
		return err
	}

	fmt.Println("debug::: 100", buf)

	end := bytes.IndexByte(buf, EOM)
	if end < 0 {
		err := ProtocolError{"Message is missing EOM byte."}
		return err
	}

	payload := string(buf[:end])
	sc.feed <- payload

	if m.cmd == SUB {
		goto WAIT_FOR_SERVER
	}

	return nil
}

func (sc *Connection) Get(key string) string {
	m := NewGetMessage(key)
	err := sc.transmit(m)
	if err != nil {
		panic(err)
	}
	return <-sc.feed
}

func (sc *Connection) Set(key, value string) string {
	m := NewSetMessage(key, value)
	err := sc.transmit(m)
	if err != nil {
		panic(err)
	}
	return <-sc.feed
}

func (sc *Connection) Delete(key string) string {
	m := NewDeleteMessage(key)
	err := sc.transmit(m)
	if err != nil {
		panic(err)
	}
	return <-sc.feed
}

func (sc *Connection) Publish(key, value string) string {
	m := NewPublishMessage(key, value)
	err := sc.transmit(m)
	if err != nil {
		panic(err)
	}
	return <-sc.feed
}

func (sc *Connection) Subscribe(key string, recv chan<- string) {
	m := NewGetMessage(key)
	go sc.transmit(m)
	for {
		recv <- <-sc.feed
	}
}

func (sc *Connection) Close() error {
	err := sc.conn.Close()
	return err
}

func NewConnection(address string, port uint) (*Connection, error) {
	fullAddress := address + ":" + strconv.FormatUint(uint64(port), 10)
	conn, err := net.Dial("tcp", fullAddress)
	sc := Connection{}
	if err == nil {
		sc.conn = conn
	}
	sc.feed = make(chan string, 2)
	return &sc, err
}

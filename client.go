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
	"bytes"
	"net"
	"strconv"
)

type ServerConnection struct {
	conn net.Conn
	feed chan string
}

func (sc *ServerConnection) transmit(m *Message) error {
	b := m.Bytes()
	_, err := (*sc).conn.Write(b)

	buf := make([]byte, 1024)

WAIT_FOR_SERVER:
	_, err = (*sc).conn.Read(buf)
	if err != nil {
		return err
	}

	end := bytes.IndexByte(buf, EOM)
	if end < 0 {
		err := ProtocolError{"Message is missing EOM byte."}
		return err
	}

	payload := string(buf[:end])
	(*sc).feed <- payload

	if (*m).cmd == SUB {
		goto WAIT_FOR_SERVER
	}

	return nil
}

func (sc *ServerConnection) Get(key string) string {
	m := NewGetMessage(key)
	err := (*sc).transmit(m)
	if err != nil {
		panic(err)
	}
	return <-(*sc).feed
}

func (sc *ServerConnection) Set(key, value string) string {
	m := NewSetMessage(key, value)
	err := (*sc).transmit(m)
	if err != nil {
		panic(err)
	}
	return <-(*sc).feed
}

func (sc *ServerConnection) Delete(key string) string {
	m := NewDeleteMessage(key)
	err := (*sc).transmit(m)
	if err != nil {
		panic(err)
	}
	return <-(*sc).feed
}

func (sc *ServerConnection) Publish(key, value string) string {
	m := NewPublishMessage(key, value)
	err := (*sc).transmit(m)
	if err != nil {
		panic(err)
	}
	return <-(*sc).feed
}

func (sc *ServerConnection) Subscribe(key string, recv chan<- string) {
	m := NewGetMessage(key)
	go (*sc).transmit(m)
	for {
		value := <-(*sc).feed
		recv<-value
	}
}

func (sc *ServerConnection) Close() error {
	err := (*sc).Close()
	return err
}

func NewServerConnection(address string, port uint) (*ServerConnection, error) {
	fullAddress := address + ":" + strconv.FormatUint(uint64(port), 10)
	conn, err := net.Dial("tcp", fullAddress)
	sc := ServerConnection{}
	if err == nil {
		sc.conn = conn
	}
	return &sc, err
}

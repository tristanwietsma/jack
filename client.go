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
	"net"
	"strconv"
)

type ServerConnection struct {
	conn net.Conn
}

func (sc *ServerConnection) Send(m *Message)  {
	// to do
}

func (sc *ServerConnection) Close() {
	// to do
}

func NewServerConnection(address string, port uint) (ServerConnection, error) {
	fullAddress := address + ":" + strconv.FormatUint(uint64(port), 10)
	conn, err := net.Dial("tcp", fullAddress)
	sc = ServerConnection{}
	if err == nil {
		sc.conn = conn
	}
	return sc, err
}


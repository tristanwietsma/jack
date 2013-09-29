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

package main

import (
	"flag"
	"fmt"
	"sync"
	"strconv"
	"github.com/tristanwietsma/jackdb"
)

var addr = flag.String("address", "0.0.0.0", "server address")
var port = flag.Uint("port", 2000, "port number")
var numClients = flag.Uint("numClients", 200, "number of concurrent clients")

func main() {

	flag.Parse()

	// build client pool
	clients := []jackdb.ServerConnections{}
	for i:=0; i<numClients; i++ {
		c, err := jackdb.NewServerConnection(addr, port)
		if err != nil {
			panic(err)
		}
		clients = append(clients, c)
	}

	// set
	wg := sync.WaitGroup()
	for i:=0; i<numClients; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			key := "key" + strconv.Itoa(idx)
			_ = clients[idx].Set(key, "val")
		}(i)
	}
	wg.Wait()

	// get

	// delete

}

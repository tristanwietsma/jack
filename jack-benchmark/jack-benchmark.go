package main

import (
	"flag"
	"fmt"
	"github.com/tristanwietsma/jack"
	"strconv"
	"sync"
	"time"
)

var address = flag.String("address", "127.0.0.1", "server address")
var port = flag.Uint("port", 2000, "port number")
var numClients = flag.Uint("clients", 100, "concurrent clients")

func main() {

	flag.Parse()

	pool := jack.NewConnectionPool(*address, *port, *numClients)

	nc := int(*numClients)
	clients := []*jack.Connection{}
	for i := 0; i < nc; i++ {
		c, _ := pool.Connect()
		clients = append(clients, c)
	}

	var wg sync.WaitGroup
	work := 1000000/nc
	startTime := time.Now().UnixNano()
	for i := 0; i < nc; i++ {
		wg.Add(1)
		go func(groupId int) {
			for j := 0; j < work; j++ {
				key := strconv.Itoa(groupId) + ":" + strconv.Itoa(j)
				clients[groupId].Set(key, "abc")
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	stopTime := time.Now().UnixNano()

	tme := float64(stopTime - startTime)/1000000000
	fmt.Printf("1000000 SETS across %d clients: %.2f seconds\n", nc, tme)
}

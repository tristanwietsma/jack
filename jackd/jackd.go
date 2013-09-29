package main

import (
	"flag"
	"github.com/tristanwietsma/jackdb"
)

var port = flag.Int("port", 2000, "port number")
var buckets = flag.Int("buckets", 1000, "number of buckets")

// Main method
func main() {
	flag.Parse()
	jackdb.StartServer(*port, *buckets)
}

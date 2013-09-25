package main

import (
	"flag"
	jackdb "github.com/tristanwietsma/jackdb"
)

var port = flag.Int("port", 2000, "tcp port number")

// Main method
func main() {
	flag.Parse()
	jackdb.StartServer(*port)
}

package jackdb

import (
	"flag"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

var (
	// DB is the server's storage container.
	DB Store

	// OnlyOnce makes sure the DB is not initialized more than once.
	OnlyOnce sync.Once
)

func initialize() {
	DB.Init()
}

// StartServer starts a server on a given port.
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

// CloseConnection terminates a client's connection.
func CloseConnection(c net.Conn) {
	log.Printf("[%s] closed connection\n", c.RemoteAddr())
	c.Close()
}

// HandleConnection manages client connections to the server
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

package jackdb

import (
	"log"
	"net"
	"strconv"
	"strings"
	"github.com/tristanwietsma/metastore"
)

func StartServer(port int, buckets int) {

	var db metastore.MetaStore
	db.Init(buckets)
	log.Printf("created storage with %d buckets\n", buckets)

	portStr := ":" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", portStr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	log.Printf("server started on port %d\n", port)

	for {
		c, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go HandleConnection(c, &db)
	}
}

func CloseConnection(c net.Conn) {
	log.Printf("[%s] closed connection\n", c.RemoteAddr())
	c.Close()
}

func HandleConnection(c net.Conn, dbase *MetaStore) {
	defer CloseConnection(c)

	fromAddr := c.RemoteAddr()
	log.Printf("[%s] new connection\n", fromAddr)

	buf := make([]byte, 1024)

	bucketIndex := dbase.GetHasher()

NEXTMESSAGE:

	_, err := c.Read(buf)
	if err != nil {
		return
	}

	msg, err := Parse(buf)
	if err != nil {
		return
	}

	i := bucketIndex(msg.key)

	switch msg.cmd {

	case GET:

		if value, ok := (*dbase).bucket[i].Get(msg.key); ok {
			value = append(value, EOM)
			_, err = c.Write(value)
		} else {
			_, err = c.Write([]byte{EOM})
		}
		if err != nil {
			return
		}
		goto NEXTMESSAGE

	case SET:

		(*dbase).bucket[i].Set(msg.key, msg.arg)
		_, err = c.Write([]byte{SUCCESS})
		if err != nil {
			return
		}
		goto NEXTMESSAGE

	case DEL:

		(*dbase).bucket[i].Delete(msg.key)
		_, err = c.Write([]byte{SUCCESS})
		if err != nil {
			return
		}
		goto NEXTMESSAGE

	case PUB:

		(*dbase).bucket[i].Publish(msg.key, msg.arg)
		_, err = c.Write([]byte{SUCCESS})
		if err != nil {
			return
		}
		goto NEXTMESSAGE

	case SUB:

		outgoing := make(chan string)
		(*dbase).bucket[i].Subscribe(msg.key, outgoing)
		for value := range outgoing {
			value = append(value, EOM)
			_, err := c.Write(value)
			if err != nil {
				close(outgoing)
				return
			}
		}

	default:
		return
	}
}

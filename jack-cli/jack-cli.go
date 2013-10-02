package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/tristanwietsma/jackdb"
	"os"
	"os/signal"
	"strings"
	"strconv"
)

var address = flag.String("address", "127.0.0.1", "server address")
var port = flag.Uint("port", 2000, "port number")
var poolsize = flag.Uint("cmax", 100, "max connections")

func readLine() []string {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		return strings.Split(scanner.Text(), " ")
	}
	return []string{}
}

func main() {

	flag.Parse()

	pool := jackdb.NewConnectionPool(*address, *port, *poolsize)

	conn, err := pool.Connect()
	if err != nil {
		panic(err)
	}

	// handle termination signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, os.Kill)
	go func() {
		<-sigs
		fmt.Println()
		os.Exit(0)
	}()

	cmdPrompt := "\033[34;1mjack " + *address + ":" + strconv.FormatUint(uint64(*port), 10)  +">\033[0m "
	syntaxError := "\033[31;1mSyntaxError: %s\033[0m\n"

	var tokens []string
	for {
		fmt.Print(cmdPrompt)
		tokens = readLine()

		// handle empty string
		if len(tokens) == 0 {
			continue
		}

		cmd := strings.ToUpper(tokens[0])

		switch cmd {
		case "GET":

			// handle syntax error
			if len(tokens) == 1 {
				fmt.Printf(syntaxError, "GET key [key ...]")
				continue
			}

			// get each key
			for _, key := range tokens[1:] {
				value := conn.Get(key)
				if len(value) > 0 {
					fmt.Printf("%s := %s\n", key, value)
				}
			}

		case "SET":

			// handle syntax error
			if len(tokens) != 3 {
				fmt.Printf(syntaxError, "SET key value")
				continue
			}

			// set key
			fmt.Println(conn.Set(tokens[1], tokens[2]))

		case "DEL":

			// handle syntax error
			if len(tokens) == 1 {
				fmt.Printf(syntaxError, "DEL key [key ...]")
				continue
			}

			// delete each key
			for _, key := range tokens[1:] {
				fmt.Println(conn.Delete(key))
			}

		case "PUB":

			// handle syntax error
			if len(tokens) != 3 {
				fmt.Printf(syntaxError, "PUB key value")
				continue
			}

			// publish key
			fmt.Println(conn.Publish(tokens[1], tokens[2]))

		case "SUB":

			// handle syntax error
			if len(tokens) == 1 {
				fmt.Printf(syntaxError, "SUB key [key ...]")
				continue
			}

			// subscribe to each key
			recv := make(chan string)
			for _, key := range tokens[1:] {
				go func(k string) {
					c, err := pool.Connect()
					if err != nil {
						panic(err)
					}
					middle := make(chan string)
					go c.Subscribe(k, middle)
					for {
						recv <- k + " := " + <-middle
					}
				}(key)
			}

			for {
				fmt.Println(<-recv)
			}

		default:
			fmt.Printf(syntaxError, cmd)
		}
	}
}

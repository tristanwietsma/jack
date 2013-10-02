package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/tristanwietsma/jackdb"
	"os"
	"os/signal"
	"strings"
)

var address = flag.String("address", "0.0.0.0", "server address")
var port = flag.Uint("port", 2000, "port number")
var poolsize = flag.Uint("cmax", 100, "max connections")

func ReadLine() []string {
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

	cmdPrompt := "\033[34;1mjack>\033[0m "
	respPrompt := "\033[33;1mjack>\033[0m "

	var tokens []string
	for {
		fmt.Print(cmdPrompt)
		tokens = ReadLine()

		// handle empty string
		if len(tokens) == 0 {
			continue
		}

		cmd := strings.ToUpper(tokens[0])

		switch cmd {
		case "GET":

			// handle syntax error
			if len(tokens) == 1 {
				fmt.Println("\033[31;1mSyntaxError: GET key [key ...]\033[0m")
				continue
			}

			// get each key
			for _, key := range tokens[1:] {
				value := conn.Get(key)
				if len(value) > 0 {
					fmt.Printf(respPrompt + "%s := %s\n", key, value)
				}
			}

		case "SET":

			// handle syntax error
			if len(tokens) != 3 {
				fmt.Println("\033[31;1mSyntaxError: SET key value\033[0m")
				continue
			}

			// set key
			fmt.Println(respPrompt + conn.Set(tokens[1], tokens[2]))

		case "DEL":

			// handle syntax error
			if len(tokens) == 1 {
				fmt.Println("\033[31;1mSyntaxError: DEL key [key ...]\033[0m")
				continue
			}

			// delete each key
			for _, key := range tokens[1:] {
				fmt.Println(respPrompt + conn.Delete(key))
			}

		case "PUB":

			// handle syntax error
			if len(tokens) != 3 {
				fmt.Println("\033[31;1mSyntaxError: PUB key value\033[0m")
				continue
			}

			// publish key
			fmt.Println(respPrompt + conn.Publish(tokens[1], tokens[2]))

		case "SUB":

			// handle syntax error
			if len(tokens) == 1 {
				fmt.Println("\033[31;1mSyntaxError: SUB key [key ...]\033[0m")
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
				fmt.Println(respPrompt + <-recv)
			}

		default:
			fmt.Printf("\033[31;1mUnknown command: %s\033[0m\n", cmd)
		}
	}
}

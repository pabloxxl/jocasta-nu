package main

import (
	"log"
	"os"
	"strconv"

	"github.com/pabloxxl/jocasta-nu/pkg/dns"
)

func parseEnv() (int, string, int, bool) {
	portEnv, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatal("PORT environment variable is not set!")
	}

	resolverIPEnv, ok := os.LookupEnv("RESOLVER_IP")
	if !ok {
		log.Fatal("RESOLVER_IP environment variable is not set!")
	}
	resolverPortEnv, ok := os.LookupEnv("RESOLVER_PORT")
	if !ok {
		log.Fatal("RESOLVER_PORT environment variable is not set!")
	}

	debugEnv, ok := os.LookupEnv("DEBUG")
	if !ok {
		debugEnv = "0"
	}

	port, error := strconv.Atoi(portEnv)
	if error != nil {
		log.Fatal("PORT environment variable is not integer!")

	}

	resolverPort, error := strconv.Atoi(resolverPortEnv)
	if error != nil {
		log.Fatal("RESOLVER_PORT environment variable is not integer!")

	}

	debugInt, error := strconv.Atoi(debugEnv)
	if error != nil {
		log.Fatal("DEBUG environment variable is not integer!")

	}

	debug := debugInt > 0

	return port, resolverIPEnv, resolverPort, debug
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	port, resolverIP, resolverPort, debug := parseEnv()
	dnsServer := dns.GetConnection(port, resolverIP, resolverPort, debug)
	dns.Listen(dnsServer)
}

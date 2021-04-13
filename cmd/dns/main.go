package main

import (
	"log"
	"os"
	"strconv"

	"github.com/pabloxxl/jocasta-nu/pkg/dns"
)

func parseEnv() (int, string, int) {
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

	port, error := strconv.Atoi(portEnv)
	if error != nil {
		log.Fatal("PORT environment variable is not integer!")

	}

	resolverPort, error := strconv.Atoi(resolverPortEnv)
	if error != nil {
		log.Fatal("RESOLVER_PORT environment variable is not integer!")

	}

	return port, resolverIPEnv, resolverPort
}

func main() {
	blockedHosts := dns.CreateAllRecordsFromDB()
	log.Printf("Read %d records from database", len(*blockedHosts))
	port, resolverIP, resolverPort := parseEnv()
	dnsServer := dns.GetConnection(port, resolverIP, resolverPort, blockedHosts)
	dns.Listen(dnsServer)
}

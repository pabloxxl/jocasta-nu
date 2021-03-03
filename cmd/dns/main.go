package main

import (
	"bufio"
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

func readBlockedFile(path string) []dns.Record {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Cannot open %v", path)
	}
	defer file.Close()

	var lines []dns.Record
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, *dns.CreateRecordBlock(scanner.Text()))
	}
	return lines
}

func main() {
	blockedHosts := readBlockedFile("/blocked.txt")
	port, resolverIP, resolverPort := parseEnv()
	dnsServer := dns.GetConnection(port, resolverIP, resolverPort, blockedHosts)
	dns.Listen(dnsServer)
}

package main

import (
	"log"

	"github.com/pabloxxl/jocasta-nu/pkg/rest"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	rest.Serve()
}

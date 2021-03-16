package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pabloxxl/jocasta-nu/pkg/db"
	"github.com/pabloxxl/jocasta-nu/pkg/dns"
)

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling /ping")
	client := db.CreateClient()
	defer db.DisconnectClient(client)

	w.Write([]byte(db.GetDatabaseNames(client)))

}
func insert(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling /insert")

	url, ok := r.URL.Query()["url"]

	if !ok || len(url[0]) < 1 {
		log.Println("Missing parameter: url")
		return
	}

	client := db.CreateClient()

	record := dns.Record{Action: dns.ActionBlock, URL: url[0]}
	db.PutRecord(client, record)
}

// Serve serve rest api
func Serve() {
	log.Printf("Listening on port %d", 8080)
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/ping", ping)
	myRouter.HandleFunc("/insert", insert)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

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

	record := dns.CreateRecordBlock(url[0])
	log.Printf("Inserting record: %v", record)
	db.PutAny(client, record)
}

func records(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling /records")

	client := db.CreateClient()

	records := db.GetAny(client, "records", "action", 0)

	for _, value := range records {
		actionInt := int(value["action"].(int32))
		record := dns.CreateRecord(value["url"].(string), actionInt)
		w.Write([]byte(fmt.Sprintf("%s %s\n", dns.ActionToString(record.Action), record.URL)))
	}
}

// Serve serve rest api
func Serve() {
	log.Printf("Listening on port %d", 8080)
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/ping", ping)
	myRouter.HandleFunc("/insert", insert)
	myRouter.HandleFunc("/records", records)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

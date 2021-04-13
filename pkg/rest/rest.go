package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pabloxxl/jocasta-nu/pkg/db"
	"github.com/pabloxxl/jocasta-nu/pkg/dns"
	"go.mongodb.org/mongo-driver/mongo"
)

func ping(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling /ping")
	client := db.CreateClient()
	defer db.DisconnectClient(client)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func stats(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handling /stats")

		// TODO maybe use json here?
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("secret_number: 42\n"))
		w.Write([]byte(fmt.Sprintf("number_of_records: %d\n", db.CountDocuments(client, "records"))))
	}
}

func insert(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handling /insert")

		url, ok := r.URL.Query()["url"]

		if !ok || len(url[0]) < 1 {
			log.Println("Missing parameter: url")
			return
		}

		action := dns.ActionBlock
		actionFromURL, ok := r.URL.Query()["action"]
		if ok && len(actionFromURL[0]) > 0 {
			action = dns.StringToAction(actionFromURL[0])
		}

		questionRecord := dns.GetOneRecordFromDB(url[0])
		if !dns.IsRecordEmpty(questionRecord) {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(fmt.Sprintf("%d: Conflict with %s", http.StatusConflict, dns.RecordToString(questionRecord))))
			return
		}

		record := dns.CreateRecord(url[0], action)
		db.PutAny(client, "records", record)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%d: OK", http.StatusOK)))
	}
}

func clear(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handling /clear")

		db.DeleteAll(client, "records")

		w.WriteHeader(http.StatusOK)
		// TODO it would be nice to print number of deleted entries in response
		w.Write([]byte(fmt.Sprintf("%d: OK", http.StatusOK)))
	}
}

func records(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handling /records")

		records := *dns.CreateAllRecordsFromDB(client)

		w.WriteHeader(http.StatusOK)
		for _, value := range records {
			w.Write([]byte(fmt.Sprintf("%s %s\n", dns.ActionToString(value.Action), value.URL)))
		}
	}
}

// Serve serve rest api
func Serve() {
	log.Printf("Listening on port %d", 8080)
	client := db.CreateClient()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/ping", ping)
	router.HandleFunc("/insert", insert(client))
	router.HandleFunc("/clear", clear(client))
	router.HandleFunc("/records", records(client))
	router.HandleFunc("/stats", stats(client))
	log.Fatal(http.ListenAndServe(":8080", router))
}

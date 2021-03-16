package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pabloxxl/jocasta-nu/pkg/db"
)

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling /ping")
	client := db.CreateClient()
	defer db.DisconnectClient(client)

	w.Write([]byte(db.GetDatabaseNames(client)))

}

// Serve serve rest api
func Serve() {
	log.Printf("Listening on port %d", 8080)
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/ping", ping)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

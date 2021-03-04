package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling /ping")
	w.Write([]byte("pong"))
}

// Serve serve rest api
func Serve() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/ping", ping)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

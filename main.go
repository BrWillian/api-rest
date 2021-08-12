// main.go
package main

import (
	"log"
	"net/http"

	routers "github.com/brwillian/api-rest/routers"
	"github.com/gorilla/mux"
)

func handleRequests() {
	myRouter := mux.NewRouter()
	myRouter.HandleFunc("/api/classificador/veicular", routers.GetClassificacao).Methods("POST")
	myRouter.HandleFunc("/api/version", routers.GetVersion)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	handleRequests()
}

// main.go
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	routers "github.com/brwillian/api-rest/routers"
)

func handleRequests() {
	myRouter := mux.NewRouter()
	myRouter.HandleFunc("/api/classificador", routers.getClassificacao).Methods("POST")
	myRouter.HandleFunc("/api/version", routers.getVersion)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	handleRequests()
}

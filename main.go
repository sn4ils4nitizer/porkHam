package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"porkHam/handlers"
)

func main() {
	// Initialize mux router, will define api endpoint
	router := mux.NewRouter()

	//Routes - takes URL, GetPage etc... defined elswhere (probably handlers.go lol)
	router.HandleFunc("/api/wiki/list", handlers.ListPages).Methods("GET")
	router.HandleFunc("/api/wiki/{name}", handlers.GetPage).Methods("GET")
	router.HandleFunc("/api/wiki/{name}", handlers.CreatePage).Methods("POST")
	router.HandleFunc("/api/wiki/{name}", handlers.DeletePage).Methods("DELETE")
	router.HandleFunc("/api/wiki/{name}", handlers.ModifyPage).Methods("PUT")

	// Start server
	// If server fails to start error is logges and program crashes,
	// listen and serve starts HTTP server on port 8080
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))

}

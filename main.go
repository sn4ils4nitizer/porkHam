package main

import (
	gorillaHandler "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"porkHam/handlers"
)

func main() {
	// Initialize mux router, will define api endpoint
	router := mux.NewRouter()

	// WIKI ROUTES
	// Methods defined in porkHam/handlers/wiki.go
	router.HandleFunc("/api/wiki2/list", handlers.ListPages).Methods("GET")
	router.HandleFunc("/api/wiki/{name}", handlers.GetPage).Methods("GET")
	router.HandleFunc("/api/wiki/{name}", handlers.CreatePage).Methods("POST")
	router.HandleFunc("/api/wiki/{name}", handlers.DeletePage).Methods("DELETE")
	router.HandleFunc("/api/wiki/{name}", handlers.ModifyPage).Methods("PUT")

	// TODO
	// WEATHER ROUTES
	// Methods defined in porkHam/handlers/weather.go
	router.HandleFunc("/api/weather/{location}", handlers.WeatherHandler).Methods("GET")

	// ORACLE ROUTES
	// Methods defined in handlers/oracle.go
	router.HandleFunc("/api/oracle/prompt", handlers.QueryHandler).Methods("POST")
	router.PathPrefix("/oracleaudio/").Handler(http.StripPrefix("/oracleaudio/", http.FileServer(http.Dir("oracleaudio"))))

	cors := gorillaHandler.CORS(
		gorillaHandler.AllowedOrigins([]string{"*"}),
		gorillaHandler.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		gorillaHandler.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	router.PathPrefix("/pages/").Handler(http.StripPrefix("/pages", http.FileServer(http.Dir("./pages"))))
	// Start server
	// If server fails to start error is logges and program crashes,
	// listen and serve starts HTTP server on port 8080
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", cors(router)))

}

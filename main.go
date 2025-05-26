package main

import (
	gorillaHandler "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"porkHam/handlers"
)

func main() {
	router := mux.NewRouter()

	// Your routes here...
	router.HandleFunc("/api/wiki2/list", handlers.ListPages).Methods("GET")
	router.HandleFunc("/api/wiki/{name}", handlers.GetPage).Methods("GET")
	router.HandleFunc("/api/wiki/{name}", handlers.CreatePage).Methods("POST")
	router.HandleFunc("/api/wiki/{name}", handlers.DeletePage).Methods("DELETE")
	router.HandleFunc("/api/wiki/{name}", handlers.ModifyPage).Methods("PUT")

	router.HandleFunc("/api/weather/{location}", handlers.WeatherHandler).Methods("GET")

	router.HandleFunc("/api/oracle/prompt", handlers.QueryHandler).Methods("POST", "OPTIONS")

	router.PathPrefix("/oracleaudio/").Handler(http.StripPrefix("/oracleaudio/", http.FileServer(http.Dir("oracleaudio"))))
	router.PathPrefix("/pages/").Handler(http.StripPrefix("/pages", http.FileServer(http.Dir("./pages"))))

	corsMiddleware := gorillaHandler.CORS(
		gorillaHandler.AllowedOrigins([]string{"*"}),
		gorillaHandler.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		gorillaHandler.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)(router) // wrap once

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/oracle/prompt" {
			router.ServeHTTP(w, r) // oracle route, no global CORS — it's handled inside QueryHandler
		} else {
			corsMiddleware.ServeHTTP(w, r) // global CORS for all others
		}
	})

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

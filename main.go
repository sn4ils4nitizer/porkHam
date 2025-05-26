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

	// Wiki routes
	router.HandleFunc("/api/wiki2/list", handlers.ListPages).Methods("GET")
	router.HandleFunc("/api/wiki/{name}", handlers.GetPage).Methods("GET")
	router.HandleFunc("/api/wiki/{name}", handlers.CreatePage).Methods("POST")
	router.HandleFunc("/api/wiki/{name}", handlers.DeletePage).Methods("DELETE")
	router.HandleFunc("/api/wiki/{name}", handlers.ModifyPage).Methods("PUT")

	// Weather route
	router.HandleFunc("/api/weather/{location}", handlers.WeatherHandler).Methods("GET")

	// Oracle route - handles its own CORS inside the handler
	router.HandleFunc("/api/oracle/prompt", handlers.QueryHandler).Methods("POST", "OPTIONS")

	// Static file handlers
	router.PathPrefix("/oracleaudio/").Handler(http.StripPrefix("/oracleaudio/", http.FileServer(http.Dir("oracleaudio"))))
	router.PathPrefix("/pages/").Handler(http.StripPrefix("/pages", http.FileServer(http.Dir("./pages"))))

	// Global CORS middleware for all routes except /api/oracle/prompt
	corsMiddleware := gorillaHandler.CORS(
		gorillaHandler.AllowedOrigins([]string{"*"}),
		gorillaHandler.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		gorillaHandler.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)(router)

	// Custom handler to skip global CORS on oracle route (because it sets its own)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/oracle/prompt" {
			log.Println("Serving oracle route without global CORS")
			router.ServeHTTP(w, r) // No global CORS, handler adds it manually
		} else {
			corsMiddleware.ServeHTTP(w, r) // Global CORS applied
		}
	})

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

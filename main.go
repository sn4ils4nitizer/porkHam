package main

import (
	gorillaHandler "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"porkHam/handlers"
)

func main() {

	var wikiPath = "/api/wiki/"
	// Router with global CORS middleware applied (for all routes except oracle)
	globalRouter := mux.NewRouter()

	// Define routes that get global CORS
	globalRouter.HandleFunc(wikiPath+"list", handlers.ListPages).Methods("GET")
	globalRouter.HandleFunc(wikiPath+"{name}", handlers.GetPage).Methods("GET")
	globalRouter.HandleFunc(wikiPath+"{name}", handlers.CreatePage).Methods("POST")
	globalRouter.HandleFunc(wikiPath+"{name}", handlers.DeletePage).Methods("DELETE")
	globalRouter.HandleFunc(wikiPath+"{name}", handlers.ModifyPage).Methods("PUT")

	globalRouter.HandleFunc("/api/weather/{location}", handlers.WeatherHandler).Methods("GET")

	globalRouter.PathPrefix("/oracleaudio/").Handler(http.StripPrefix("/oracleaudio/", http.FileServer(http.Dir("oracleaudio"))))
	globalRouter.PathPrefix("/pages/").Handler(http.StripPrefix("/pages", http.FileServer(http.Dir("./pages"))))

	// Apply global CORS middleware only here
	corsMiddleware := gorillaHandler.CORS(
		gorillaHandler.AllowedOrigins([]string{"*"}),
		gorillaHandler.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		gorillaHandler.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)(globalRouter)

	// Separate router for oracle route with no global CORS middleware
	oracleRouter := mux.NewRouter()
	oracleRouter.HandleFunc("/api/oracle/prompt", handlers.QueryHandler).Methods("POST", "OPTIONS")

	// Now use a top-level router to delegate based on path prefix
	mainRouter := mux.NewRouter()
	mainRouter.PathPrefix("/api/oracle/prompt").Handler(oracleRouter)
	mainRouter.PathPrefix("/").Handler(corsMiddleware) // everything else gets global CORS

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mainRouter))
}

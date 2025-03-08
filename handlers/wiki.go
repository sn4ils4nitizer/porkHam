package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	//"os"
	"porkHam/utils"

	"github.com/gorilla/mux"
)

// Method definitions
// GetPage - extracts name from the URL, uses utils.ReadFile(name) to get the page
// shows message is file does not exist, if exists sets Content-Type: text/html, return content
func GetPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)  // Vars extract URL params, part of Mux
	name := vars["name"] // name is set to vars - which is extracted from r

	//if there is no file
	content, err := utils.ReadFile(name) // tries to read file
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(content))
}

func CreatePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err = utils.WriteFile(name, string(body))
	if err != nil {
		log.Printf("Error saving page %s: %v", name, err) // Log the actual error
		http.Error(w, "Failed to save page", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Page %s saved successfully!", name)
}

// Delete page
func DeletePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	err := utils.DeleteFile(name)
	if err != nil {
		http.Error(w, "Failed to delete page", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Page %s deleted successfully!", name)
}

func ModifyPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	//Read body of the request - this is the new content for the page
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	//Add new contecnt to file
	err = utils.WriteFile(name, string(body))
	if err != nil {
		http.Error(w, "Failed to modify page", http.StatusInternalServerError)
		return
	}

	//Success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Page %s modified successfully!", name)
}

func ListPages(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for /api/wiki/list") // Debugging line
	files, err := ioutil.ReadDir("pages")
	if err != nil {
		http.Error(w, "Could not read pages directory", http.StatusInternalServerError)
		return
	}
	var pageNames []string
	for _, file := range files {
		if file.IsDir() {
			continue // skip directories in /pages
		}
		pageNames = append(pageNames, file.Name())
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pageNames)
}

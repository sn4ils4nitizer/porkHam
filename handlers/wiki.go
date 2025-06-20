package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"os"
	//"path/filepath"
	"porkHam/utils"

	"github.com/gorilla/mux"
)

// Method definitions
// GetPage - extracts name from the URL, uses utils.ReadFile(name) to get the page
// shows message is file does not exist, if exists sets Content-Type: text/html, return content
func GetPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rpath := vars["path"]

	log.Println("Requested url: " + rpath)

	filePath := filepath.Join(rpath)
	log.Println("Full path to file: ", filePath)

	content, err := utils.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		log.Println("Error reading file: ", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(content))
}

// Second version of CreatePage. This version creates directories too!
func CreatePage(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	path := vars["path"]
	name := vars["name"]

	// Clean up spaces from path and name
	path = strings.TrimSpace(path)
	name = strings.TrimSpace(name)

	fullpath := filepath.Join("./pages", path, name+".html")

	log.Println("Saving file to: ", fullpath)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err = utils.WriteFile(fullpath, string(body))
	if err != nil {
		log.Printf("Error saving page %s: %v", fullpath, err)
		http.Error(w, "Failed to save the page", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Page %s saved successfully!", name)
}

// Delete page
func DeletePage(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	path := vars["path"]
	rpath := "./pages/" + path

	log.Println("Requested deletion of: ", rpath)

	err := os.Remove(rpath)
	if err != nil {
		http.Error(w, "Failed to delete page: "+rpath, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Page %s deleted successfully. \n", rpath)

	err = utils.DeleteEmptyDirs(filepath.Dir(rpath), "./pages")
	if err != nil {
		log.Println("Error while cleaning up empty directories: ", err)
	}
}
func ListPages(w http.ResponseWriter, r *http.Request) {

	log.Println("Received request for /api/wiki/list")

	tree, err := utils.BuildTree("pages")
	if err != nil {
		log.Printf("Unable to build tree: %v", err)
		http.Error(w, "Could not build tree", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tree)
	if err != nil {
		http.Error(w, "Failed to encode file tree", http.StatusInternalServerError)
	}
}

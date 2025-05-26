package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func QueryHandler(w http.ResponseWriter, r *http.Request) {

	// CORS preflight check
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Regular CORS headers for actual request
	w.Header().Set("Access-Control-Allow-Origin", "*")

	prompt, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	cmd := exec.Command("./oracleQuery.pl", string(prompt))
	if err := cmd.Run(); err != nil {
		http.Error(w, "Oracle Query Command Failed."+err.Error(), http.StatusInternalServerError)
		return
	}

	audiodata, err := os.ReadFile("output.wav")
	if err != nil {
		http.Error(w, "Failed to fetch audio:"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "audio/wave")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(audiodata)))
	w.WriteHeader(http.StatusOK)

	w.Write(audiodata)
	log.Print(w, "Oracle query sent seccessfully.")
}

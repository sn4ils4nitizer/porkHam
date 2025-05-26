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

package handlers

import (
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

func QueryHandler(w http.ResponseWriter, r *http.Request) {

	prompt, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	cmd := exec.Command("./oracleQuery.pl ", string(prompt))
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, "Oracle Query Command Failed."+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(output)
	log.Print(w, "Oracle query sent seccessfully.")
}

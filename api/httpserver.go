package api

import (
	"encoding/json"
	"net/http"
	"project3/pkg/common"
	"project3/pkg/groundstation"
)

// StartAPIServer Starting the HTTP API server, this is a scalable function
func StartAPIServer() {
	http.HandleFunc("/vessels", handleVessels)
	common.Logger.Println("API server started at :12345")
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		return
	}
}

func handleVessels(w http.ResponseWriter, r *http.Request) {
	messages, err := groundstation.LoadFromDatabase()
	if err != nil {
		http.Error(w, "Failed to load data", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(messages)
}

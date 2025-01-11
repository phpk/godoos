package proxy

import (
	"encoding/json"
	"godo/progress"
	"log"
	"net/http"
)

func StartFrpcHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if err := progress.StartCmd("frpc"); err != nil {
		log.Printf("Failed to start frpc service: %v", err)
		http.Error(w, "Failed to start frpc service", http.StatusInternalServerError)
		return
	}

	response := NewResponse("0", "frpc service started", nil)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func StopFrpcHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if err := progress.StopCmd("frpc"); err != nil {
		log.Printf("Failed to stop frpc service: %v", err)
		http.Error(w, "Failed to stop frpc service", http.StatusInternalServerError)
		return
	}

	response := NewResponse("0", "frpc service stoped", nil)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func RestartFrpcHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if err := progress.RestartCmd("frpc"); err != nil {
		log.Printf("Failed to restarted frpc service: %v", err)
		http.Error(w, "Failed to restarted frpc service", http.StatusInternalServerError)
		return
	}

	response := NewResponse("0", "frpc service restarted", nil)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

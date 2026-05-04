package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/quangho/yoga-be/internal/adapter/firestoredb"
)

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]any{"error": err.Error()})
}

func handleRepoError(w http.ResponseWriter, err error) {
	if errors.Is(err, firestoredb.ErrNotFound) {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeError(w, http.StatusBadRequest, err)
}

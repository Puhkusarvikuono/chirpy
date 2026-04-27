package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(403)
		w.Write([]byte("Forbidden"))
		return
	}
	err := cfg.db.Reset(r.Context())
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error resetting database"))
	}
	w.Write([]byte("Database reset\n"))
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

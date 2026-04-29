package main

import (
	"log"
	"net/http"

	"github.com/puhkusarvikuono/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Unauthorized: %v\n", err)
		w.WriteHeader(401)
	}

	refreshToken, err := cfg.db.GetUserFromRefreshToken(r.Context(), bearerToken)
	if err != nil {
		log.Printf("Unauthorized: %v\n", err)
		w.WriteHeader(401)
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken.Token)
	if err != nil {
		log.Printf("Token revoke unsuccessful: %v\n", err)
		w.WriteHeader(500)
	}

	log.Printf("Token revoked successfully")
	w.WriteHeader(204)
}

package main

import (
	"log"
	"net/http"
	"time"

	"github.com/puhkusarvikuono/chirpy/internal/auth"
)

type rToken struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
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

	if refreshToken.RevokedAt.Valid {
		log.Println("Unauthorized, refresh token revoked")
		w.WriteHeader(401)
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		log.Println("Unauthorized, refresh token expired")
		w.WriteHeader(401)
	}

	expiresAt := 1 * time.Hour

	token, err := auth.MakeJWT(refreshToken.UserID, cfg.secret, expiresAt)
	if err != nil {
		log.Printf("Error creating JWT token: %v\n", err)
		w.WriteHeader(401)
	}

	response := rToken{
		Token: token,
	}

	respondWithJSON(w, 200, response)
}

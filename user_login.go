package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"github.com/puhkusarvikuono/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
		ExpiresInSeconds	int	`json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if params.ExpiresInSeconds == 0 {
		params.ExpiresInSeconds = 3600
	}

	user, err := cfg.db.GetUser(r.Context(), params.Email)

	if err != nil {
		log.Printf("Incorrect email")
		w.WriteHeader(401)
		return
	}

	// check for hashed password match

	ok, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)

	if !ok || err != nil {
		log.Printf("Incorrect password")
		w.WriteHeader(401)
		return
	}
	
	ExpiresAt := time.Duration(params.ExpiresInSeconds) * time.Second

	if err != nil {
		log.Printf("Invalid expires at %v\n", err)
		w.WriteHeader(500)
		return
	}
	token, err := auth.MakeJWT(user.ID, cfg.secret, ExpiresAt)

	if err != nil {
		log.Printf("Error creating token %v\n", err)
		w.WriteHeader(500)
		return
	}
	
	dbUser := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: token,
	}

	respondWithJSON(w, 200, dbUser)

}



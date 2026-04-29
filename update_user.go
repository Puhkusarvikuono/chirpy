package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/puhkusarvikuono/chirpy/internal/auth"
	"github.com/puhkusarvikuono/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error finding user token %v\n", err)
		w.WriteHeader(401)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		log.Printf("Not authorized: %v\n", err)
		w.WriteHeader(401)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		w.WriteHeader(500)
		return
	}

	err = cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		log.Printf("Error updating user: %s", err)
		w.WriteHeader(500)
		return
	}

	user, err := cfg.db.GetUser(r.Context(), params.Email)
	if err != nil {
		log.Printf("Error getting user: %s", err)
		w.WriteHeader(500)
		return
	}

	dbUser := User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		Token:       token,
		IsChirpyRed: user.IsChirpyRed,
	}

	respondWithJSON(w, 200, dbUser)
}

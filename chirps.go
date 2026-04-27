package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"
	"github.com/google/uuid"
	"github.com/puhkusarvikuono/chirpy/internal/database"
	"time"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body     	string    `json:"body"`
	UserID		uuid.UUID	`json:"user_id"`
}

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		ID	uuid.UUID	`json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if utf8.RuneCountInString(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	cleanMsg := profanityCheck(params.Body, badWords)
	if err != nil {
		log.Printf("Error parsing user id: %s", err)
		w.WriteHeader(500)
		return
	}


	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleanMsg,
		UserID: params.ID,
	})

	if err != nil {
		log.Printf("Error creating user: %s", err)
		w.WriteHeader(500)
		return
	}
	
	dbChirp := Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	}



	respondWithJSON(w, 201, dbChirp)
}



func profanityCheck(msg string, badWords map[string]struct{}) string {
	words := strings.Split(msg, " ")

	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}

	cleaned := strings.Join(words, " ")
	return cleaned
}

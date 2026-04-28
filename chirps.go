package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/puhkusarvikuono/chirpy/internal/database"
	"github.com/puhkusarvikuono/chirpy/internal/auth"
	"time"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
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

	
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleanMsg,
		UserID: userID,
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
		UserID: userID,
	}

	chirp := databaseChirpToChirp(dbChirp)

	respondWithJSON(w, 201, chirp)
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

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		log.Printf("Error getting chirps: %s", err)
		w.WriteHeader(500)
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirp := databaseChirpToChirp(dbChirp)
		chirps = append(chirps, chirp)
	}

	respondWithJSON(w, 200, chirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	pathname := r.PathValue("chirpID")
	target, err := uuid.Parse(pathname)
	if err != nil {
		log.Printf("Error parsing chirp id: %s", err)
		w.WriteHeader(500)
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), target)
	if err != nil {
		log.Printf("%s\n", target)
		log.Printf("Chirp not found: %s", err)
		w.WriteHeader(404)
	}

	chirp := databaseChirpToChirp(dbChirp)
	respondWithJSON(w, 200, chirp)
}

func databaseChirpToChirp(dbChirp database.Chirp) Chirp {
	return Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		UserID:    dbChirp.UserID,
		Body:      dbChirp.Body,
	}
}

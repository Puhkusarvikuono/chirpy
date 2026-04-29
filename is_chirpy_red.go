package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/puhkusarvikuono/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerIsChirpyRed(w http.ResponseWriter, r *http.Request) {
	APIKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		log.Printf("Unauthorized: %v\n", err)
		w.WriteHeader(401)
	}

	if APIKey != cfg.polkaKey {
		log.Printf("Unauthorized")
		w.WriteHeader(401)
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	target, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		log.Printf("Error parsing user id: %v\n", err)
		w.WriteHeader(500)
		return
	}

	err = cfg.db.UpgradeRedUser(r.Context(), target)
	if err != nil {
		log.Printf("User not found")
		w.WriteHeader(404)
		return
	}

	log.Printf("User upgraded successfully.")
	w.WriteHeader(204)
}

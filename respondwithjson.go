package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"slices"
	"time"
	"github.com/google/uuid"
	"github.com/odilmode/http/internal/auth"
	"github.com/odilmode/http/internal/database"

)


type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error encoding response: %w", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	return respondWithJSON(w, code, map[string]string{"error": msg})
}

func (cfg *apiConfig) handleChirps(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type requestBody struct {
		Body string `json:"body"`
	}
	type responseBody struct {
		CleanedBody string `json:"cleaned_body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, " Couldn't validate JWT")
		return
	}
	var params requestBody
	if err = json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode request")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedText := wordreplace(params.Body)
	now := time.Now().UTC()
	chirpParams := database.CreateChirpParams{
    		ID:        uuid.New(),
    		CreatedAt: now,
    		UpdatedAt: now,
    		Body:      cleanedText,    // sanitized string
   		UserID:    userID,  // from request
	}
	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		fmt.Println("CreateChirp DB error:", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	})
}

func wordreplace(sentence string) string {
	badWords := []string{"kerfuffle","sharbert","fornax"}

	words := strings.Split(sentence, " ")
	for i, word := range words {
		if slices.Contains(badWords, strings.ToLower(word)) {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

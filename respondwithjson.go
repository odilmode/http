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

// Chirp represents a microblog post (chirp) returned in responses
// @Description A chirp created by a user
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
// requestBody represents the JSON request body for creating a chirp
type requestBody struct {
	// Body is the text content of the chirp
	// max length: 140 characters
	Body string `json:"body"`
}
// responseBody represents the JSON response body after creating a chirp
type responseBody struct {
	// CleanedBody is the sanitized chirp text with bad words replaced
	CleanedBody string `json:"cleaned_body"`
}
// handleChirps creates a new chirp
// @Summary      Create a new chirp
// @Description  Authenticated endpoint to create a chirp with max length 140 characters. Filters bad words.
// @Tags         chirps
// @Accept       json
// @Produce      json
// @Param        chirp  body requestBody true "Chirp body"
// @Success      201  {object}  Chirp
// @Failure      400  {object}  ErrorResponse "Invalid author_id"
// @Failure      401  {object}  map[string]string  "Unauthorized - missing or invalid JWT"
// @Failure      500  {object}  map[string]string  "Internal server error - failed to create chirp"
// @Security     BearerAuth
// @Router       /api/chirps [post]
func (cfg *apiConfig) handleChirps(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
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
// ErrorResponse represents an error response message
// swagger:model ErrorResponse
type ErrorResponse struct {
    // Error message describing what went wrong
    Error string `json:"error"`
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

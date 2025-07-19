package main


import (
	"log"
	"encoding/json"
	"net/http"
	"github.com/google/uuid"
)
// handleGetChirp returns a single chirp by ID.
// @Summary Get a chirp
// @Description Retrieve a chirp by its ID
// @Tags Chirps
// @Accept json
// @Produce json
// @Param chirpID path string true "Chirp ID"
// @Success 200 {object} Chirp
// @Failure 400 {object} map[string]string "Invalid chirp ID"
// @Failure 404 {object} map[string]string "Chirp not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/chirps/{chirpID} [get]
func (cfg *apiConfig) handleGetChirp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	chirpID := r.PathValue("chirpID")
	id, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}
	chirp, err := cfg.dbQueries.GetChirp(ctx, id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}
	responseChirp := Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:	chirp.Body,
		UserID: chirp.UserID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(responseChirp); err != nil {
		log.Println("Failed to encode chirp:", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to encode chirp")
	}
}

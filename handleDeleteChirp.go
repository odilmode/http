package main

import (
	"net/http"
	"github.com/odilmode/http/internal/auth"
	"github.com/google/uuid"
)
// handleDeleteChirp deletes a chirp by its ID if the requester is the author.
// @Summary Delete a chirp
// @Description Delete a chirp if the authenticated user is the author
// @Tags Chirps
// @Accept json
// @Produce json
// @Param chirpID path string true "Chirp ID"
// @Param Authorization header string true "Bearer JWT token"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "Invalid chirp ID"
// @Failure 401 {object} map[string]string "Unauthorized or missing token"
// @Failure 403 {object} map[string]string "Forbidden: not the author"
// @Failure 404 {object} map[string]string "Chirp not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/chirps/{chirpID} [delete]
func (cfg *apiConfig) handleDeleteChirp(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or Invalid Authorization header")
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

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

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "The user is not the author")
		return
	}
	err = cfg.dbQueries.DeleteChirp(ctx, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

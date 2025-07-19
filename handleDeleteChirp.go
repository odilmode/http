package main

import (
	"net/http"
	"github.com/odilmode/http/internal/auth"
	"github.com/google/uuid"
)

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

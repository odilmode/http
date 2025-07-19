package main

import (
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"github.com/odilmode/http/internal/auth"
)

func (cfg *apiConfig) handleWebhooks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var requestBody struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get api key")
		return
	}

	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized user")
		return
	}

	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode request")
		return
	}

	if requestBody.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(requestBody.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	err = cfg.dbQueries.UpgradeUserToChirpyRed(ctx, userID)
	if err != nil {
		if err.Error() == "user not found" {
			respondWithError(w, http.StatusNotFound, "Failed to upgrade user")
			return
		}
		log.Printf("UpgradeUserToChirpyRed error: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/*func (cfg *apiConfig) handleWebhooks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var requestBody struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode request")
		return
	}

	if requestBody.Event != "user.upgraded" {
		// Ignore other events with 204 no content
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(requestBody.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	err = cfg.dbQueries.UpgradeUserToChirpyRed(ctx, userID)
	if err != nil {
		if err.Error() == "user not found" {
			respondWithError(w, http.StatusNotFound, "Failed to upgrade user")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Success, no content response
	w.WriteHeader(http.StatusNoContent)
}
*/

package main
import (
	"log"
	"github.com/odilmode/http/internal/database"
	"github.com/google/uuid"
	"encoding/json"
	"net/http"
	"sort"
)
// handleGetChirps godoc
// @Summary      Get Chirps
// @Description  Retrieve chirps, optionally filtered by author_id and sorted by creation time
// @Tags         chirps
// @Accept       json
// @Produce      json
// @Param        author_id  query     string  false  "Filter chirps by author UUID"
// @Param        sort       query     string  false  "Sort order: asc (default) or desc"
// @Success      200        {array}   Chirp
// @Failure      400        {object}  ErrorResponse "Invalid author_id"
// @Failure      500        {object}  ErrorResponse "Failed to fetch chirps or encode response"
// @Router       /api/chirps [get]
func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s := r.URL.Query().Get("author_id")
	var chirps []database.Chirp
	var err error
	if s != "" {
		authorID, err := uuid.Parse(s)
		if err !=  nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author_id")
			return
		}
		chirps, err = cfg.dbQueries.GetChirpsByAuthor(ctx, authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to fetch chirps")
			return
		}
	} else {
		chirps, err = cfg.dbQueries.GetAllChirps(ctx)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to fetch chirps")
			return
		}
	}
	sortOrder := r.URL.Query().Get("sort")
	if sortOrder != "desc" {
		// Default to ascending
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		})
	} else {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[j].CreatedAt.Before(chirps[i].CreatedAt)
		})
	}

	responseChirps := []Chirp{}
	for _, c := range chirps {
		responseChirps = append(responseChirps, Chirp{
			ID: c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body: c.Body,
			UserID: c.UserID,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(responseChirps); err != nil {
		log.Println("Failed to encode chirps:", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to encode chirps")
	}
}

package main
import (
	"github.com/odilmode/http/internal/auth"
	"net/http"
)

func (cfg *apiConfig) handleRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or Invalid Authorization header")
		return
	}

	ctx := r.Context()

	err = cfg.dbQueries.RevokeRefreshToken(ctx, refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to revoke token")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

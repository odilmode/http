package main
import (
	"github.com/odilmode/http/internal/auth"
	"net/http"
)
// handleRevoke godoc
// @Summary      Revoke Refresh Token
// @Description  Revokes a given refresh token, effectively logging out the user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer refresh token"
// @Success      204  "No Content"
// @Failure      401  {object}  ErrorResponse "Missing or invalid authorization header"
// @Failure      500  {object}  ErrorResponse "Failed to revoke token"
// @Router       /api/revoke [post]
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

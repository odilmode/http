package main
import (
	"github.com/odilmode/http/internal/auth"
	"net/http"
	"time"
)
// handleRefresh godoc
// @Summary      Refresh Access Token
// @Description  Generates a new access JWT given a valid refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer refresh token"
// @Success      200  {object}  map[string]string "New access token"
// @Failure      401  {object}  ErrorResponse "Missing, invalid or expired refresh token"
// @Failure      500  {object}  ErrorResponse "Failed to create access token"
// @Router       /api/refresh [post]
func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, r *http.Request) {
	refreshtoken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or Invalid Authorization header")
		return
	}


	ctx := r.Context()
	user, err := cfg.dbQueries.GetUserFromRefreshToken(ctx, refreshtoken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired refresh token")
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access token")
		return
	}

	response := map[string]string{
		"token": accessToken,
	}
	respondWithJSON(w, http.StatusOK, response)
}

package main

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/odilmode/http/internal/database"
	"github.com/odilmode/http/internal/auth"
)
// LoginRequest represents the login credentials
// swagger:model LoginRequest
type LoginRequest struct {
    Password string `json:"password"`
    Email    string `json:"email"`
}
// response represents response by server
// swagger: model response
type response struct {
	User
	Token string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

// handleLogin godoc
// @Summary      User Login
// @Description  Authenticates user and returns JWT access and refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      LoginRequest  true "User email and password"
// @Success      200          {object}  response
// @Failure      401          {object}  ErrorResponse "Incorrect email or password"
// @Failure      500          {object}  ErrorResponse "Internal server error"
// @Router       /api/login [post]
func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := LoginRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	expirationTime := time.Hour
	

	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.jwtSecret,
		expirationTime,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token")
		return
	}

	if err := cfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't save refresh token")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token: accessToken,
		RefreshToken: refreshToken,
	})
}

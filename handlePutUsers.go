package main

import (
	"net/http"
	"github.com/odilmode/http/internal/auth"
	"github.com/odilmode/http/internal/database"
	"encoding/json"
)

type RequestBody struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type ResponseBody struct {
	User
}
// handlePutUsers godoc
// @Summary      Update User Info
// @Description  Updates authenticated user's email and password
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        user body RequestBody true "Updated user email and password"
// @Success      200  {object}  ResponseBody
// @Failure      400  {object}  ErrorResponse "Invalid request body"
// @Failure      401  {object}  ErrorResponse "Unauthorized or invalid token"
// @Failure      500  {object}  ErrorResponse "Internal server error"
// @Router       /api/users [put]
func (cfg *apiConfig) handlePutUsers(w http.ResponseWriter, r *http.Request) {
	accessToken, err :=  auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or Invalid Authorization header")
		return
	}
	
	userID, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}


	decoder := json.NewDecoder(r.Body)
	params := RequestBody{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request Body")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}
	ctx := r.Context()
	updatedUser, err := cfg.dbQueries.UpdateUser(ctx, database.UpdateUserParams{
		ID: userID,
		Email: params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user params")
		return
	}
	respondWithJSON(w, http.StatusOK, ResponseBody{
		User : User{
			ID:          updatedUser.ID,
			CreatedAt:   updatedUser.CreatedAt,
			UpdatedAt:   updatedUser.UpdatedAt,
			Email:       updatedUser.Email,
			IsChirpyRed: updatedUser.IsChirpyRed,
		},
	})
}

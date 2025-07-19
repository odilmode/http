package main
import(
	"log"
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/odilmode/http/internal/auth"
	"github.com/odilmode/http/internal/database"
	"time"
)

type User struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
	Password string `json:"-"`
	IsChirpyRed bool `json:"is_chirpy_red"`
}


type createUserRequest struct {
	Password string `json:"password"`
	Email string `json:"email"`
}
// handleCreateUsers creates a new user in the system.
// @Summary Create a new user
// @Description Register a new user with email and password
// @Tags Users
// @Accept json
// @Produce json
// @Param user body createUserRequest true "User credentials"
// @Success 201 {object} User
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users [post]
func (cfg *apiConfig) handleCreateUsers(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := createUserRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		http.Error(w, "Could not decode request body", http.StatusBadRequest)
		return
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}
	user, err := cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		log.Printf("Error calling to database: %s", err)
		http.Error(w, "Failed to create user in database", http.StatusInternalServerError)
		return
	}
	mainUser := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:	user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}
	jsonData, err := json.Marshal(mainUser)
	if err != nil {
		log.Printf("Error converting to json: %s", err)
		http.Error(w, "Failed to marshal response data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(jsonData)
	if err != nil {
		log.Printf("Error writing response: %s", err)
		return
	}

}	

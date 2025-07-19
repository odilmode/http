package main
import "net/http"
import "log"
// resetMetrics godoc
// @Summary      Reset users and file server hits
// @Description  Deletes all users from the database and resets hit counter. Only accessible in development environment.
// @Tags         admin
// @Produce      plain
// @Success      200  {string}  string  "All users deleted and hits reset to 0"
// @Failure      403  {string}  string  "Forbidden: This endpoint is only accessible in development environment"
// @Failure      500  {string}  string  "Failed to delete all users"
// @Router       /admin/reset [post]
func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		http.Error(w, "Forbidden: Thid endpoint is only accessible in development environment", http.StatusForbidden)
		return
	}

	// deleting users from database
	err := cfg.dbQueries.DeleteAllUsers(r.Context())
	if err != nil {
		log.Printf("Error deleting all users: %s", err)
		http.Error(w, "Failed to delete all users", http.StatusInternalServerError)
		return
	}

	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All users deleted and hits reset to 0"))
}

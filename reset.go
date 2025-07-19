package main
import "net/http"
import "log"

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

package main
import (
	"net/http"
	"fmt"
)

// handleReadiness godoc
// @Summary      Readiness check endpoint
// @Description  Returns "OK" if the server is ready to handle requests
// @Tags         health
// @Produce      plain
// @Success      200  {string}  string  "OK"
// @Router       /api/healthz [get]
func handleReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	bodyText := "OK"
	_, err := w.Write([]byte(bodyText))
	if err != nil {
		fmt.Printf("Error writing response: %v\n", err)
	}
}


package main
import (
	"fmt"
	"net/http"
)
// middlewareMetricsInc is an HTTP middleware that increments a counter for hits.
// It wraps another http.Handler and increments the fileserverHits counter on each request.
//
// Note: Middleware functions generally don't get swagger docs as endpoints themselves.

// handleMetrics godoc
// @Summary      Show Chirpy usage metrics
// @Description  Returns an HTML page showing how many times Chirpy has been visited
// @Tags         admin
// @Produce      html
// @Success      200  {string}  string  "HTML content with visit count"
// @Router       /admin/metrics [get]
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())))
}


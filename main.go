package main
import _ "github.com/lib/pq"
import (
	"log"
	"net/http"
	"sync/atomic"
	"os"
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/odilmode/http/internal/database"
)


type apiConfig struct {
	fileserverHits		atomic.Int32
	writeHandler		string
	dbQueries		*database.Queries
	Platform		string
	jwtSecret		string
	polkaKey		string
}


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error connecting to database: %s\n", err)
	}

	dbQueries := database.New(db)
	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir(filepathRoot))

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not secret")
	}
	polka := os.Getenv("POLKA_KEY")
	if polka == "" {
		log.Fatal("Polka_Key is not set")
	}
	platform := os.Getenv("PLATFORM")
	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries:	dbQueries,
		Platform:	platform,
		jwtSecret:	jwtSecret,
		polkaKey:	polka,
	}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fs)))
	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetMetrics)
	mux.HandleFunc("POST /api/users", apiCfg.handleCreateUsers)
	mux.HandleFunc("POST /api/chirps", apiCfg.handleChirps)
	mux.HandleFunc("GET /api/chirps", apiCfg.handleGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handleGetChirp)
	mux.HandleFunc("POST /api/login", apiCfg.handleLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handleRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handleRevoke)
	mux.HandleFunc("PUT /api/users", apiCfg.handlePutUsers)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handleDeleteChirp)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handleWebhooks)


	server := &http.Server{
		Addr:		":" + port,
		Handler:	mux,
	}
	
	// Use http.RedirectHandler() function to create a handler which 307
	// redirects all requests it receives to http://example.org.
	//rh := http.RedirectHandler("http://example.org", 307)

	// Next we use mux.Handle() function to register this with our new
	// servemux, so it acts as the handler for all incoming requests with URL
	// path /foo

	
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	// Then we create a new server and start listening for incoming requests
	// with the http.ListenAndServe() function, passing in our servemux for it to
	// match requests against as the second argument.
	log.Fatal(server.ListenAndServe())
}


package main

import (
	"log"
	"database/sql"
	"net/http"
	"os"
	"sync/atomic"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"github.com/puhkusarvikuono/chirpy/internal/database"
)


type apiConfig struct {
	fileserverHits atomic.Int32
	db	*database.Queries
	platform string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)

	if err != nil {
		log.Fatalf("Error: %v", err)
	} 

	const port = "8080"
	apiCfg := apiConfig{
		db: dbQueries,
		platform: platform,
	}
	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirps)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

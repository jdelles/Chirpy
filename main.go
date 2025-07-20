package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"Chirpy/handlers"
	"Chirpy/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
    apiCfg := &handlers.ApiConfig{}
    godotenv.Load()
    
    dbURL := os.Getenv("DB_URL")
    jwtSecret := os.Getenv("JWT_SECRET")
    polkaKey := os.Getenv("POLKA_KEY")
	platform := os.Getenv("PLATFORM")
    
    if dbURL == "" {
        log.Fatal("DB_URL must be set")
    }
    if jwtSecret == "" {
        log.Fatal("JWT_SECRET must be set")
    }
    if polkaKey == "" {
        log.Fatal("POLKA_KEY must be set")
    }
	if platform == "" {
        log.Fatal("PLATFORM must be set")
    }
    
    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    
    dbQueries := database.New(db)
    apiCfg.DbQueries = dbQueries
    apiCfg.JwtSecret = jwtSecret
    apiCfg.PolkaKey = polkaKey
	apiCfg.Platform = platform

	filepathRoot := "."
	port := "8080"

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.MiddlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlers.HandlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.HandlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.HandlerReset)
	mux.HandleFunc("POST /api/users", apiCfg.HandlerCreateUser)
	mux.HandleFunc("POST /api/chirps", apiCfg.HandlerCreateChirps)
	mux.HandleFunc("GET /api/chirps", apiCfg.HandlerReadChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.HandlerReadChirpsByID)
	mux.HandleFunc("POST /api/login", apiCfg.HandlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.HandlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.HandlerRevoke)
	mux.HandleFunc("PUT /api/users", apiCfg.HandlerUpdateUsers)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.HandlerDeleteChirps)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.HandlerPolkaWebhook)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

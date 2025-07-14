package main

import (
	"encoding/json"
	"fmt"
	"log"
    "os"
	"net/http"
    "database/sql"
	"strings"
	"sync/atomic"

    "github.com/joho/godotenv"
    _ "github.com/lib/pq"
    "Chirpy/internal/database"
)

func main() {
    apiCfg := &apiConfig{}
    godotenv.Load()
    dbURL := os.Getenv("DB_URL")
    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    dbQueries := database.New(db)
    apiCfg.dbQueries = dbQueries

	filepathRoot := "."
	port := "8080"
	
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
        Body string `json:"body"`
    }
    params := parameters{}

    type returnVals struct {
        Cleaned_body string `json:"cleaned_body"`
    }
    
    decoder := json.NewDecoder(r.Body)
    
    err := decoder.Decode(&params)
    if err != nil {
        log.Printf("Error decoding parameters %s", err)
        respondWithError(w, http.StatusBadRequest, "Something went wrong")
        return
    }
    
    if len(params.Body) > 140 {
        respondWithError(w, http.StatusBadRequest, "Chirp is too long")
        return
    }

    replacement := "****"
    badWords := []string{"kerfuffle", "sharbert", "fornax"}
    splitString := strings.Split(params.Body, " ")
    cleanedWords := make([]string, len(splitString))
    for i, word := range splitString {
        found := false
        for _, bad := range badWords {
            if strings.ToLower(word) == bad {
                found = true
                break
            }
        }
        if found {
            cleanedWords[i] = replacement
        } else {
            cleanedWords[i] = word
        }
    }
    result := strings.Join(cleanedWords, " ")

    respondWithJSON(w, http.StatusOK, returnVals{
        Cleaned_body: result,
    })
}
type apiConfig struct {
	fileserverHits atomic.Int32
    dbQueries      *database.Queries
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := cfg.fileserverHits.Load()
	response := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, hits)
	w.Write([]byte(response))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
    type errorResponse struct {
        Error string `json:"error"`
    }
    respondWithJSON(w, code, errorResponse{Error: msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    dat, err := json.Marshal(payload)
    if err != nil {
        log.Printf("Error marshaling JSON: %s", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    w.Write(dat)
}

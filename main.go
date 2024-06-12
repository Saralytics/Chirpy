// godotenv.Load()
package main

import (
	"chirpy/m/internal/database"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// by default, godotenv will look for a file named .env in the current directory

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	jwtKey         string
	polkaApiKey    string
	refresh_token  string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	apiKey := os.Getenv("POLKA_APIKEY")

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatalf("Error initializing database : %v", err)
	}

	apiCfg := &apiConfig{
		fileserverHits: 0,
		DB:             db,
		jwtKey:         jwtSecret,
		polkaApiKey:    apiKey,
	}

	mux := http.NewServeMux()
	fileServerHandler := http.FileServer(http.Dir("."))
	wrappedFileServerHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", fileServerHandler))
	mux.Handle("/app/", wrappedFileServerHandler)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("/api/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsGet)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerChirpGetByID)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerChirpDelete)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	mux.HandleFunc("POST /api/login", apiCfg.handlerUsersLogin)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUsersUpdate)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerTokenRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerTokenRevoke)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerPolkaWebhook)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Println("Listening...")
	log.Fatal(server.ListenAndServe())

}

package main

import (
	"chirpy/m/internal/database"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func main() {
	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatalf("Error initializing database : %v", err)
	}

	apiCfg := &apiConfig{
		fileserverHits: 0,
		DB:             db,
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
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	mux.HandleFunc("POST /api/login", apiCfg.handlerUsersLogin)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Println("Listening...")
	log.Fatal(server.ListenAndServe())

}

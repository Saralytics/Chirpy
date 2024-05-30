package main

import (
	"chirpy/m/internal/database"
	// "chirpy/m/internal/utils"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

const metricsTemplate = `<html>

<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>

</html>
`

func main() {
	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatalf("Error initializing database : %v", err)
	}

	handler := &Handler{
		DB: db,
	}

	apiCfg := &apiConfig{}

	mux := http.NewServeMux()
	fileServerHandler := http.FileServer(http.Dir("."))
	wrappedFileServerHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", fileServerHandler))
	mux.Handle("/app/", wrappedFileServerHandler)
	mux.HandleFunc("GET /api/healthz", checkHealthHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("/api/reset", apiCfg.resetHandler)
	mux.HandleFunc("POST /api/chirps", handler.validationHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Println("Listening...")
	log.Fatal(server.ListenAndServe())

}

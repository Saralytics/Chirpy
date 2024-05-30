package main

import (
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
	apiCfg := &apiConfig{}

	mux := http.NewServeMux()
	fileServerHandler := http.FileServer(http.Dir("."))
	wrappedFileServerHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", fileServerHandler))
	mux.Handle("/app/", wrappedFileServerHandler)
	mux.HandleFunc("GET /api/healthz", checkHealthHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("/api/reset", apiCfg.resetHandler)
	mux.HandleFunc("POST /api/chirps", validationHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Println("Listening...")
	log.Fatal(server.ListenAndServe())

}

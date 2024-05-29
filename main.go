package main

import (
	"fmt"
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

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {

	newHandler := func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(newHandler)
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	metrics := fmt.Sprintf(metricsTemplate, cfg.fileserverHits)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(metrics))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits = 0
}

func main() {
	apiCfg := &apiConfig{}

	checkHealthHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}

	mux := http.NewServeMux()
	fileServerHandler := http.FileServer(http.Dir("."))
	wrappedFileServerHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", fileServerHandler))
	mux.Handle("/app/", wrappedFileServerHandler)
	mux.HandleFunc("GET /api/healthz", checkHealthHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("/api/reset", apiCfg.resetHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Println("Listening...")
	log.Fatal(server.ListenAndServe())

}

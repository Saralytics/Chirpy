package main

import (
	"chirpy/m/internal/database"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var profane_words = [3]string{"kerfuffle", "sharbert", "fornax"}

type errorResp struct {
	Error string `json:"error"`
}

type Handler struct {
	DB *database.DB
}

func checkHealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

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

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	errorRespBody := errorResp{
		Error: msg,
	}
	resp, _ := json.Marshal(errorRespBody)

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	respBody := payload
	resp, _ := json.Marshal(respBody)

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)

}

func (h *Handler) validationHandler(w http.ResponseWriter, r *http.Request) {

	type payload struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := payload{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	// validation

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	// if the body contains the bad words

	words := strings.Split(params.Body, " ")
	for _, profane_word := range profane_words {
		for i, word := range words {
			if strings.ToLower(word) == strings.ToLower(profane_word) {
				// replace the word with "****"
				words[i] = "****"
			}
		}
	}
	body := strings.Join(words, " ")
	newChirp, err := h.DB.CreateChirp(body)

	respondWithJSON(w, 201, newChirp)
}

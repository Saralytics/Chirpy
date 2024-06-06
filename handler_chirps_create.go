package main

import (
	"chirpy/m/internal/auth"
	"chirpy/m/internal/database"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {

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

	// authenticate the user and get back a user id
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Not authorized to perform this action")
	}
	userIDStr, err := auth.ValidateJWT(token, cfg.jwtKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	userIDInt, err := strconv.Atoi(userIDStr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	body, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	cleaned := cleanBody(body)

	newChirp, err := cfg.DB.CreateChirp(cleaned, userIDInt)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create chirp")
	}

	respondWithJSON(w, http.StatusCreated, database.Chirp{
		ID:       newChirp.ID,
		Body:     newChirp.Body,
		AuthorID: newChirp.AuthorID})
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {

		return "", errors.New("chirp is too long")
	}
	return body, nil
}

func cleanBody(body string) string {
	// if the body contains the bad words
	words := strings.Split(body, " ")
	for _, profane_word := range profane_words {
		for i, word := range words {
			if strings.ToLower(word) == strings.ToLower(profane_word) {
				// replace the word with "****"
				words[i] = "****"
			}
		}
	}
	cleanedBody := strings.Join(words, " ")
	return cleanedBody
}

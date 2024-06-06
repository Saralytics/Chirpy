package main

import (
	"chirpy/m/internal/auth"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) handlerChirpDelete(w http.ResponseWriter, r *http.Request) {

	// get the token from header
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}
	// authenticate user -> get user id
	userIDStr, err := auth.ValidateJWT(token, cfg.jwtKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	userIDInt, err := strconv.Atoi(userIDStr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// get the chirp id from request
	chirpID, err := strconv.Atoi(r.PathValue("chirpID"))

	// call the delete
	err = cfg.DB.DeleteChirpByID(userIDInt, chirpID)
	if err != nil {
		respondWithError(w, 403, "")
	}
	respondWithJSON(w, 204, nil)

}

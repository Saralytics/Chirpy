package main

import (
	"chirpy/m/internal/auth"
	"net/http"
)

func (cfg *apiConfig) handlerTokenRefresh(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "the token is invalid")
	}
	tokenStr, err := cfg.DB.RefreshJWT(token, cfg.jwtKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	payload := struct {
		Token string `json:"token"`
	}{
		Token: tokenStr,
	}

	respondWithJSON(w, http.StatusOK, payload)

}

func (cfg *apiConfig) handlerTokenRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "the token is invalid")
	}

	err = cfg.DB.RevokeToken(token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	respondWithJSON(w, 204, nil)

}

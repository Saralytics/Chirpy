package main

import (
	"chirpy/m/internal/auth"
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	// parse the request - get the email and password from the request
	type payload struct {
		PASSWORD string `json:"password"`
		EMAIL    string `json:"email"`
	}
	var params payload
	decoder := json.NewDecoder(r.Body).Decode(&params)

	err := decoder
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "The request is invalid")
	}

	user, refresh_token, err := cfg.DB.LoginUser(params.EMAIL, params.PASSWORD)
	createdToken, _ := auth.CreateJWT(user.ID, cfg.jwtKey)

	response := struct {
		ID           int    `json:"id"`
		Email        string `json:"email"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
		IsChirpyRed  bool   `json:"is_chirpy_red"`
	}{
		ID:           user.ID,
		Email:        user.Email,
		Token:        createdToken,
		RefreshToken: refresh_token.Token,
		IsChirpyRed:  user.IsChirpyRed,
	}

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}
	cfg.refresh_token = refresh_token.Token
	respondWithJSON(w, http.StatusOK, response)

}

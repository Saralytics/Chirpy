package main

import (
	"chirpy/m/internal/auth"
	"encoding/json"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	// parse the request
	type payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var params payload
	decoder := json.NewDecoder(r.Body).Decode(&params)
	err := decoder
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "The request is invalid")
		return
	}

	// extract the token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	// validate the token
	// get user id
	userIDString, err := auth.ValidateJWT(token, cfg.jwtKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "failed to validate jwt")
		return
	}

	userIDInt, err := strconv.Atoi(userIDString)
	if err != nil {

		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// hash password
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(params.Password), 0)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}
	hashedPwdString := string(pwdHash)
	// call update user method of the database

	user, err := cfg.DB.UpdateUser(userIDInt, params.Email, hashedPwdString)

	// return updated user and 200
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
		return
	}

	response := struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}{
		ID:    user.ID,
		Email: user.Email,
	}
	respondWithJSON(w, http.StatusOK, response)
}

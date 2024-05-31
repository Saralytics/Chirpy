package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	// parse the request
	type payload struct {
		PASSWORD string `json:"password"`
		EMAIL    string `json:"email"`
	}
	var params payload
	decoder := json.NewDecoder(r.Body).Decode(&params)

	err := decoder
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	pwdHash, err := bcrypt.GenerateFromPassword([]byte(params.PASSWORD), 0)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}
	hashedPwdString := string(pwdHash)

	// save it to the database

	data, err := cfg.DB.CreateUser(params.EMAIL, hashedPwdString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating the user")
	}

	response := struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}{
		ID:    data.ID,
		Email: data.Email,
	}
	respondWithJSON(w, http.StatusCreated, response)

}

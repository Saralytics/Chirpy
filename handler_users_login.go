package main

import (
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

	data, err := cfg.DB.LoginUser(params.EMAIL, params.PASSWORD)
	response := struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}{
		ID:    data.ID,
		Email: data.Email,
	}

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}
	respondWithJSON(w, http.StatusOK, response)

}

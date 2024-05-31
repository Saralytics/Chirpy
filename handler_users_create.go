package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	// parse the request
	type payload struct {
		EMAIL string `json:"email"`
	}
	var params payload
	decoder := json.NewDecoder(r.Body).Decode(&params)

	err := decoder
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	// save it to the database
	// fmt.Println("Saved to database... %s", params.EMAIL)
	resp, err := cfg.DB.CreateUser(params.EMAIL)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating the user")
	}
	respondWithJSON(w, http.StatusCreated, resp)

}

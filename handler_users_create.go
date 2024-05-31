package main

import (
	"encoding/json"
	"fmt"
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
		respondWithError(w, http.StatusInternalServerError, "Error parsing the request")
	}

	// save it to the database
	fmt.Println("Saved to database... %s", params.EMAIL)

	// return 201

	// return id and email
}

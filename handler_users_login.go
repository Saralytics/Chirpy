package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	// parse the request - get the email and password from the request
	type payload struct {
		PASSWORD string `json:"password"`
		EMAIL    string `json:"email"`
		// ExpiresInSeconds *int32 `json:"expires_in_seconds,omitempty"`
	}
	var params payload
	decoder := json.NewDecoder(r.Body).Decode(&params)

	err := decoder
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "The request is invalid")
	}

	// if params.ExpiresInSeconds == nil {
	// 	defaultExpiration := int32(1 * 60 * 60) //access token expires after 1 hour
	// 	params.ExpiresInSeconds = &defaultExpiration
	// }

	data, err := cfg.DB.LoginUser(params.EMAIL, params.PASSWORD)
	createdToken, _ := cfg.createToken(w, data.ID)

	response := struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
		Token string `json:"token"`
	}{
		ID:    data.ID,
		Email: data.Email,
		Token: createdToken,
	}

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}
	respondWithJSON(w, http.StatusOK, response)

}

func (cfg *apiConfig) createToken(w http.ResponseWriter, userID int) (string, error) {
	key := []byte(cfg.jwtKey)

	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(int32(1*60*60)) * time.Second)),
		Subject:   fmt.Sprintf("%d", userID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(key)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return "", err
	}
	return signedToken, nil
}

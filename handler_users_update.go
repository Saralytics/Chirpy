package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ErrNoAuthHeaderIncluded -
var ErrNoAuthHeaderIncluded = errors.New("not auth header included in request")

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
	token, err := GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	// validate the token
	// get user id
	userIDString, err := validateJWT(token, cfg.jwtKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
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

func validateJWT(tokenString, jwtSecret string) (string, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(jwtSecret), nil },
	)
	if err != nil {
		return "", err
	}

	userIDString, err := token.Claims.GetSubject()
	fmt.Print(userIDString)
	if err != nil {
		return "", err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return "", err
	}
	if issuer != string("chirpy") {
		return "", errors.New("invalid issuer")
	}

	return userIDString, nil
}

// GetBearerToken -
func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}

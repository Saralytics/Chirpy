package database

import (
	"chirpy/m/internal/auth"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

var ErrTokenExpiredOrNotValid = errors.New("token expired or not valid")
var ErrTokenNotFound = errors.New("token not found")

func (db *DB) GenerateRefreshToken(userID int, expiry time.Time) (RefreshToken, error) {
	dbStructure, err := db.LoadDB()
	if err != nil {
		return RefreshToken{}, err
	}

	id := len(dbStructure.RefreshTokens) + 1

	tokenLen := 32
	arr := make([]byte, tokenLen)
	_, err = rand.Read(arr) // Use the underscore to ignore the number of bytes read
	if err != nil {
		return RefreshToken{}, err
	}

	// Convert random bytes to a hexadecimal string
	tokenStr := hex.EncodeToString(arr)
	newToken := RefreshToken{
		ID:     id,
		UserID: userID,
		Token:  tokenStr,
		Expiry: expiry,
	}
	dbStructure.RefreshTokens[id] = newToken
	err = db.writeDB(dbStructure)
	if err != nil {
		return RefreshToken{}, err
	}

	return newToken, nil
}

// func check refresh token validity
func (db *DB) RefreshJWT(token string, jwtSecret string) (string, error) {
	curDB, err := db.LoadDB()
	if err != nil {
		return "", err
	}
	for _, currentToken := range curDB.RefreshTokens {
		if currentToken.Token == token && currentToken.Expiry.After(time.Now()) {
			new_jwt, err := auth.CreateJWT(currentToken.UserID, jwtSecret)
			if err != nil {
				return "", err
			}

			return new_jwt, nil
		}
	}
	return "", ErrTokenExpiredOrNotValid
}

// func delete the expired token
func (db *DB) RevokeToken(token string) error {
	curDB, err := db.LoadDB()
	if err != nil {
		return err
	}
	tokenFound := false

	for id, currentToken := range curDB.RefreshTokens {
		if currentToken.Token == token {
			tokenFound = true
			delete(curDB.RefreshTokens, id)
			break
		}
	}
	if !tokenFound {
		return ErrTokenNotFound
	}

	err = db.writeDB(curDB)
	if err != nil {
		return err
	}
	return nil
}

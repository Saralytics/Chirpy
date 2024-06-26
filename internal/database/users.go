package database

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (db *DB) CreateUser(email string, passwordHash string) (User, error) {

	// load the db
	curDB, err := db.LoadDB()
	if err != nil {
		return User{}, err
	}

	// check if the email already exists
	for _, user := range curDB.Users {
		if email == user.Email {
			return User{}, errors.New("the email already exists")
		}
	}

	newID := len(curDB.Users) + 1

	newUser := User{
		ID:           newID,
		Email:        email,
		PasswordHash: passwordHash,
		IsChirpyRed:  false,
	}

	curDB.Users[newID] = newUser

	err = db.writeDB(curDB)
	if err != nil {
		return User{}, err
	}
	return newUser, nil
}

func (db *DB) LoginUser(email, password string) (User, RefreshToken, error) {
	// search for the email in the db
	curDB, err := db.LoadDB()
	if err != nil {
		return User{}, RefreshToken{}, err
	}
	for _, user := range curDB.Users {
		if user.Email == email {
			err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

			if err != nil {
				return User{}, RefreshToken{}, errors.New("Unauthorized")
			}

			// generate a refresh token if the use can be authenticated
			expiry := time.Now().AddDate(0, 0, 60)
			refresh_token, err := db.GenerateRefreshToken(user.ID, expiry)
			return user, refresh_token, nil
		}
	}
	return User{}, RefreshToken{}, errors.New("user not found")
}

func (db *DB) UpdateUser(id int, email string, passwordHash string) (User, error) {
	// load the current db
	// find the user id
	// update the info
	// write it back
	curDB, err := db.LoadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range curDB.Users {
		if user.ID == id {
			user.Email = email
			user.PasswordHash = passwordHash
		}
		curDB.Users[id] = user
		err = db.writeDB(curDB)
		if err != nil {
			return User{}, err
		}
		return user, nil

	}
	return User{}, errors.New("user not found")

}

func (db *DB) UpgradeUser(id int) error {
	// load the current db
	// find the user id
	// update the info
	// write it back
	curDB, err := db.LoadDB()
	if err != nil {
		return err
	}

	for _, user := range curDB.Users {
		if user.ID == id {
			user.IsChirpyRed = true

		}
		curDB.Users[id] = user
		err = db.writeDB(curDB)
		if err != nil {
			return err
		}
		return nil

	}
	return errors.New("user not found")

}

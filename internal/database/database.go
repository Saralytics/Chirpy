package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	if err := db.ensureDB(); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.LoadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	newChirp := Chirp{
		ID:   id,
		Body: body,
	}

	dbStructure.Chirps[id] = newChirp
	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return newChirp, nil

}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.LoadDB()
	if err != nil {
		return nil, err
	}
	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil

}

func (db *DB) GetChirpByID(id int) (Chirp, error) {
	dbStructure, err := db.LoadDB()
	if err != nil {
		return Chirp{}, err
	}

	for _, chirp := range dbStructure.Chirps {
		if chirp.ID == id {
			return chirp, nil
		}

	}

	return Chirp{}, errors.New("this chirp is not found")

}

func (db *DB) CreateDB() error {
	dbStructure := DBStructure{
		Chirps: map[int]Chirp{},
		Users:  map[int]User{},
	}

	err := db.writeDB(dbStructure)
	return err
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {

	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.CreateDB()
	}
	return err
}

func (db *DB) LoadDB() (DBStructure, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	// read the database file into memory json -> struct
	var dbStructure DBStructure
	data, err := os.ReadFile(db.path)
	if err != nil {
		if os.IsNotExist(err) {
			db.ensureDB()
		}
	}
	err = json.Unmarshal(data, &dbStructure)
	if err != nil {
		return dbStructure, err
	}
	return dbStructure, nil

}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	//load the entire db
	// update the data
	// Marshal back to json
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	err = os.WriteFile(db.path, data, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) CreateUser(email string) (User, error) {

	// load the db
	curDB, err := db.LoadDB()
	if err != nil {
		return User{}, err
	}

	newID := len(curDB.Users) + 1

	newUser := User{
		ID:    newID,
		EMAIL: email,
	}

	curDB.Users[newID] = newUser

	err = db.writeDB(curDB)
	if err != nil {
		return User{}, err
	}
	return newUser, nil
}

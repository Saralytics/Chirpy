package database

import (
	"encoding/json"
	"os"
	"sync"
)

type DB struct {
	path   string
	mux    *sync.RWMutex
	chirps map[int]Chirp
	nextID int
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	newChirp := Chirp{
		ID:   db.nextID,
		Body: body,
	}
	db.chirps[db.nextID] = newChirp
	db.nextID++

	if err := db.writeDB(); err != nil {
		return Chirp{}, err
	}
	return newChirp, nil
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := &DB{
		path:   path,
		mux:    &sync.RWMutex{},
		chirps: make(map[int]Chirp),
		nextID: 1,
	}

	if err := db.ensureDB(); err != nil {
		return nil, err
	}

	DBStructure, err := db.loadDB()
	if err != nil {
		return nil, err

	}

	db.chirps = DBStructure.Chirps
	for id := range db.chirps {
		if id >= db.nextID {
			db.nextID = id + 1
		}
	}
	return db, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	db.mux.Lock()
	defer db.mux.Unlock()

	if _, err := os.Stat(db.path); os.IsNotExist(err) {
		initialDB := DBStructure{
			Chirps: make(map[int]Chirp),
		}
		data, err := json.Marshal(initialDB)
		if err != nil {
			return err
		}

		err = os.WriteFile(db.path, data, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
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
func (db *DB) writeDB() error {
	//load the entire db
	// update the data
	// Marshal back to json
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStructure := DBStructure{
		Chirps: db.chirps,
	}

	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	err = os.WriteFile(db.path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

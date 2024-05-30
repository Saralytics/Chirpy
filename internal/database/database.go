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

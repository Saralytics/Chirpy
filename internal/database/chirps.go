package database

import "errors"

func (db *DB) CreateChirp(body string, author_id int) (Chirp, error) {
	dbStructure, err := db.LoadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	newChirp := Chirp{
		ID:       id,
		Body:     body,
		AuthorID: author_id,
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

func (db *DB) DeleteChirpByID(userID int, chirpID int) error {
	// find the author id of the chirp
	curDB, err := db.LoadDB()
	if err != nil {
		return err
	}

	chirpToDelete := curDB.Chirps[chirpID]
	authorID := chirpToDelete.AuthorID

	if authorID != userID {
		return errors.New("user cannot perform this action on this resource")
	}
	delete(curDB.Chirps, chirpID)
	return nil
}

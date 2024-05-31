package main

import (
	"chirpy/m/internal/database"
	"net/http"
	"sort"
	"strconv"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting from database")

	}

	chirps := []database.Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, database.Chirp{
			ID:   dbChirp.ID,
			Body: dbChirp.Body,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJSON(w, http.StatusOK, chirps)

}

func (cfg *apiConfig) handlerChirpGetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "The chirp ID is invalid")
	}
	dbChirp, err := cfg.DB.GetChirpByID(id)

	if err != nil {
		if err.Error() == "this chirp is not found" {
			respondWithError(w, http.StatusNotFound, "")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Unable to convert to JSON")
		return
	}

	respondWithJSON(w, http.StatusOK, dbChirp)

}

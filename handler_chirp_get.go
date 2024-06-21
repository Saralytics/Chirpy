package main

import (
	"chirpy/m/internal/database"
	"net/http"
	"sort"
	"strconv"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {

	author_id_str := r.URL.Query().Get("author_id")
	sort_order_str := r.URL.Query().Get("sort")
	if len(sort_order_str) == 0 {
		sort_order_str = "asc"
	}

	if len(author_id_str) != 0 {
		// get chirp by author
		author_id_int, err := strconv.Atoi(author_id_str)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "The author ID is invalid")
		}
		chirps, err := getChirpsByAuthor(cfg, author_id_int, sort_order_str)
		if err != nil {
			respondWithError(w, http.StatusNotFound, err.Error())
		}
		respondWithJSON(w, http.StatusOK, chirps)

	}

	chirps, err := getAllChirps(cfg, sort_order_str)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "")
	}
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

func getAllChirps(cfg *apiConfig, sort_order string) ([]database.Chirp, error) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		return nil, err
	}

	chirps := []database.Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, database.Chirp{
			ID:       dbChirp.ID,
			Body:     dbChirp.Body,
			AuthorID: dbChirp.AuthorID,
		})
	}

	if sort_order == "asc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID < chirps[j].ID
		})
	} else {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[j].ID < chirps[i].ID
		})
	}
	return chirps, nil
}

func getChirpsByAuthor(cfg *apiConfig, author_id int, sort_order string) ([]database.Chirp, error) {
	dbChirps, err := cfg.DB.GetChirpsByAuthor(author_id)
	if err != nil {
		return nil, err
	}

	chirps := []database.Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, database.Chirp{
			ID:       dbChirp.ID,
			Body:     dbChirp.Body,
			AuthorID: dbChirp.AuthorID,
		})
	}

	if sort_order == "asc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID < chirps[j].ID
		})
	} else {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[j].ID < chirps[i].ID
		})
	}

	return chirps, nil
}

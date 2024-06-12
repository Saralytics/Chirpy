package main

import (
	"encoding/json"
	"net/http"
)

const UpgradedStatus = "user.upgraded"

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}
	var params payload
	decoder := json.NewDecoder(r.Body).Decode(&params)

	err := decoder
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	if params.Event != UpgradedStatus {
		respondWithJSON(w, 204, nil)
	}

	if params.Event == UpgradedStatus {
		err = cfg.DB.UpgradeUser(params.Data.UserID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "user not found")
		}
		respondWithJSON(w, 204, nil)
	}
}

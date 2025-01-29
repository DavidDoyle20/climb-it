package main

import (
	"climb_it/internal/auth"
	"climb_it/internal/database"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerHabitsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name       string    `json"name"`
		Start_Date time.Time `json:"start_date"`
		End_Date   time.Time `json:"end_date"`
	}
	token, err := auth.GetAuthorizationHeader("Bearer", r.Header)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusUnauthorized, "Couldn't get bearer token from header")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Failed to validate jwt")
		return
	}

	var params parameters
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode json body")
		return
	}

	habit, err := cfg.DB.CreateHabitForUser(r.Context(), database.CreateHabitForUserParams{
		ID:        uuid.New().String(),
		Name:      params.Name,
		UserID:    userID,
		StartDate: params.Start_Date,
		EndDate:   params.End_Date,
	})
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create habit")
		return
	}
	respondWithJSON(w, http.StatusOK, habit)
}

func (cfg *apiConfig) handlerHabitsDelete(w http.ResponseWriter, r *http.Request) {
	habitIDStr := strings.TrimPrefix(r.URL.Path, "/v1/habits/")

	token, err := auth.GetAuthorizationHeader("Bearer", r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get bearer token")
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate jwt")
		return
	}

	habit, err := cfg.DB.GetHabit(r.Context(), habitIDStr)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find a habit with the given id")
		return
	}

	if userID != habit.UserID {
		respondWithError(w, http.StatusUnauthorized, "UserId and habit userID do not match")
		return
	}

	err = cfg.DB.RemoveHabit(r.Context(), habit.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't remove habit")
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}

func (cfg *apiConfig) handlerHabitsGet(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetAuthorizationHeader("Bearer", r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get bearer token")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate jwt")
		return
	}

	habits, err := cfg.DB.GetUserHabits(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user habits")
		return
	}

	respondWithJSON(w, http.StatusOK, habits)
}

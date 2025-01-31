package main

import (
	"climb_it/internal/auth"
	"climb_it/internal/database"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type paramaters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	var params paramaters
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameters")
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:             uuid.New().String(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Name:           params.Name,
		Email:          params.Email,
		HashedPassword: params.Password,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		ID            string    `json:"id"`
		Created_at    time.Time `json:"created_at"`
		Updated_at    time.Time `json:"updated_at"`
		Email         string    `json:"email"`
		Name          string    `json:"name"`
		Token         string    `json:"token"`
		Refresh_token string    `json:"refresh_token"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't read json body")
		return
	}

	user, err := cfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		log.Println(params.Email)
		respondWithError(w, http.StatusNotFound, "Couldn't find a user with that email")
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't make JWT")
		return
	}

	refresh_token, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to make refresh token")
		return
	}
	// Remove previous token
	cfg.DB.RevokeRefreshTokenFromUser(r.Context(), user.ID)
	_, err = cfg.DB.AssignRefreshTokenToUser(r.Context(), database.AssignRefreshTokenToUserParams{
		Token:  refresh_token,
		UserID: user.ID,
	})
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't assign refresh token")
		return
	}

	userAndToken := response{
		ID:            user.ID,
		Created_at:    user.CreatedAt,
		Updated_at:    user.UpdatedAt,
		Email:         user.Email,
		Name:          user.Name,
		Token:         token,
		Refresh_token: refresh_token,
	}

	respondWithJSON(w, http.StatusOK, userAndToken)
}

func (cfg *apiConfig) handlerUsersLogout(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetAuthorizationHeader("Bearer", r.Header)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusUnauthorized, "No bearer token found in header")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate jwt")
		return
	}

	err = cfg.DB.RevokeRefreshTokenFromUser(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't revoke access token. Is the user logged in")
		return
	}
	respondWithJSON(w, http.StatusOK, nil)
}

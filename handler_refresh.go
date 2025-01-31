package main

import (
	"climb_it/internal/auth"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type resp struct {
		Token string `json:"token"`
	}
	// gets bearer token "jwt" from header
	token, err := auth.GetAuthorizationHeader("Bearer", r.Header)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No token found in bearer header")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate jwt")
		return
	}

	refreshToken, err := cfg.DB.GetRefreshTokenFromUser(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldnt find a token for that user")
		return
	}

	log.Println(refreshToken)

	_, err = cfg.DB.CheckAndFetchRefreshToken(r.Context(), refreshToken.Token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "refresh token is expired or invalid")
		return
	}

	access_token, err := auth.MakeJWT(userID, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not make new jwt")
		return
	}

	respondWithJSON(w, http.StatusOK, resp{
		Token: access_token,
	})
}

package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jpforbes5151/Servertests/internal/auth"
	"github.com/jpforbes5151/Servertests/internal/database"
)

type User struct {
	Email    string `json:"email"`
	ID       int    `json:"id"`
	Password string `json:"-"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	user, err := cfg.DB.CreateUser(params.Email, hashedPassword)
	if err != nil {
		if errors.Is(err, database.ErrAlreadyExists) {
			respondWithError(w, http.StatusConflict, "User already exists")
			return
		}

		respondWithError(w, http.StatusInternalServerError, "Couldn't Create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		Email: user.Email,
		ID:    user.ID,
	})
}

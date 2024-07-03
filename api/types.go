package api

import "github.com/AndriiPets/chirpy/database"

type Api struct {
	database *database.DB
	fileserverHits int
	jwtSecret string
}

type returnError struct {
	Error string `json:"error"`
}

type returnClean struct {
	Clean string `json:"cleaned_body"`
}

type messageChirp struct {
	Body string `json:"body"`
}


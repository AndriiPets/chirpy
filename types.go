package main

type apiConfig struct {
	fileserverHits int
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


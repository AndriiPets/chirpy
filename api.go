package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AndriiPets/chirpy/utils"
)

func handleReady(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK!"))
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

func validateChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	decoder := json.NewDecoder(r.Body)
	message := messageChirp{}
	err := decoder.Decode(&message)
	
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, 500, "Something went wrong!")

		return
	}

	if len(message.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long!")

		return
	}

	resp := returnClean{
		Clean: utils.CleanInput(message.Body),
	}
	respondWithJSON(w, 200, resp)
}


func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)

	resp := returnError{
		Error: msg,
	}

	ret, _ := json.Marshal(resp)
	w.Write(ret)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	msg, _ := json.Marshal(payload)
	w.Write(msg)
}
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/AndriiPets/chirpy/database"
	"github.com/AndriiPets/chirpy/utils"
)

func handleReady(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK!"))
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

func getChirps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	
	db, erro := database.NewDB("database.json")
	if erro != nil {
		log.Fatalf("Unable to create a database :%s", erro)
	}
	chirps, _ := db.GetChirps()
	respondWithJSON(w, 201, chirps)
}

func getChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	
	db, erro := database.NewDB("database.json")
	if erro != nil {
		log.Fatalf("Unable to create a database :%s", erro)
	}
	chirp, err := db.GetChirp(r.PathValue("id"))
	if err != nil {
		respondWithError(w, 404, err.Error())
	}
	respondWithJSON(w, 200, chirp)
}

func postChirp(w http.ResponseWriter, r *http.Request) {
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

	db, erro := database.NewDB("database.json")
	if erro != nil {
		fmt.Printf("Unable to create a database :%s", erro)
	}
	chirp, errr := db.CreateChirp(utils.CleanInput(message.Body))

	if errr != nil {
		fmt.Printf("unable to create a chirp :%s\n", errr)
	}

	respondWithJSON(w, 200, chirp)

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
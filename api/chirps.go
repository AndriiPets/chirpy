package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/AndriiPets/chirpy/utils"
)

func (api *Api) GetChirps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	chirps, _ := api.database.GetChirps()
	respondWithJSON(w, 201, chirps)
}

func (api *Api) GetChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	chirp, err := api.database.GetChirp(r.PathValue("id"))
	if err != nil {
		respondWithError(w, 404, err.Error())
	}
	respondWithJSON(w, 200, chirp)
}


func (api *Api) PostChirp(w http.ResponseWriter, r *http.Request) {
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

	chirp, errr := api.database.CreateChirp(utils.CleanInput(message.Body))

	if errr != nil {
		fmt.Printf("unable to create a chirp :%s\n", errr)
	}

	respondWithJSON(w, 200, chirp)

}
package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/AndriiPets/chirpy/database"
	"github.com/AndriiPets/chirpy/utils"
	"golang.org/x/crypto/bcrypt"
)

type Api struct {
	database *database.DB
}

func NewApi() *Api {
	db, erro := database.NewDB("database.json")
	if erro != nil {
		log.Fatalf("Unable to create a database :%s", erro)
	}

	return &Api{
		database: db,
	}
}

func (api *Api) HandleReady(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK!"))
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

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

func (api *Api) PostUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	user := database.User{}
	err := decoder.Decode(&user)

	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, 500, "Something went wrong!")

		return
	}

	msg := validateUser(user)
	if msg != nil {
		respondWithError(w, 401, msg.Error())
	}

	//check for dplicate email
	users, _ := api.database.GetUsers()
	for _, usr := range users {
		if usr.Email == user.Email {
			respondWithError(w, 401, "User with this email already exist")
			return
		}
	}

	//hash password
	hashPass, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)

	user, errr := api.database.CreateUser(user.Email, string(hashPass))

	if errr != nil {
		fmt.Printf("unable to create a user :%s\n", errr)
	}

	userNoPass := database.UserNoPassword{
		Id: user.Id,
		Email: user.Email,
	}

	respondWithJSON(w, 201, userNoPass)
}

func (api *Api) LoginUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	user := database.User{}
	err := decoder.Decode(&user)

	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, 500, "Something went wrong!")

		return
	}

	msg := validateUser(user)
	if msg != nil {
		respondWithError(w, 401, msg.Error())
	}

	userNoPass := database.UserNoPassword{
		Id: user.Id,
		Email: user.Email,
	}

	users, _ := api.database.GetUsers()
	for _, usr := range users {
		if usr.Email == user.Email {
			ok := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(user.Password))
			if ok != nil {
				respondWithError(w, 401, "Wrong password!")
				return
			}
			respondWithJSON(w, 200, userNoPass)
			return
		}
		respondWithError(w, 404, "Email does not exist!")
		return
	}
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

func validateUser(user database.User) error {
	if user.Email == "" {
		return fmt.Errorf("You need a email to create a user!")
	}
	if user.Password == "" {
		return fmt.Errorf("You need a password to create a user!")
	}
	return nil
}
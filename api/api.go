package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/AndriiPets/chirpy/database"
)



func NewApi(jwt string) *Api {
	db, erro := database.NewDB("database.json")
	if erro != nil {
		log.Fatalf("Unable to create a database :%s", erro)
	}

	return &Api{
		database: db,
		jwtSecret: jwt,
	}
}

func (api *Api) HandleReady(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK!"))
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
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
package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/AndriiPets/chirpy/database"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type CustomClaims struct {
	Foo string `json:"foo"`
	jwt.RegisteredClaims
}

func (api *Api) PostUser(w http.ResponseWriter, r *http.Request) {
	user, err := decodeAndValidateUser(r)
	if err != nil {
		respondWithError(w, 401, err.Error())
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
	user, err := decodeAndValidateUser(r)
	if err != nil {
		respondWithError(w, 401, err.Error())
	}

	users, _ := api.database.GetUsers()
	for _, usr := range users {
		if usr.Email == user.Email {

			ok := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(user.Password))
			if ok != nil {
				respondWithError(w, 401, "Wrong password!")
				return
			}

			//create JWTtoken
			claims := CustomClaims{
				"bar",
				jwt.RegisteredClaims{
					Issuer: "chirpy",
					IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
					ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour * 24)),
					Subject: strconv.Itoa(usr.Id),
					},
			}
			token := jwt.NewWithClaims(
				jwt.SigningMethodHS256,
				claims,
			)

			ss, jwterr := token.SignedString([]byte(api.jwtSecret))
			if jwterr != nil {
				log.Printf("Failed to create a jwt token: %s", jwterr)
				respondWithError(w, 500, "Something went wrong!")

				return
			}

			userNoPass := database.UserNoPassword{
				Id: usr.Id,
				Email: usr.Email,
				JWTtoken: ss,
			}

			respondWithJSON(w, 200, userNoPass)
			return
		}
		respondWithError(w, 404, "Email does not exist!")
		return
	}
}

func (api *Api) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Extract JWT token
	head := r.Header.Get("Authorization")
	if head == "" {
		respondWithError(w, 401, "Cannot get the JWT token!")
		return
	}

	cleanHead, ok := strings.CutPrefix(head, "Bearer ")
	if !ok {
		respondWithError(w, 401, "Cannot get the JWT token!")
		return
	}

	token, err := jwt.ParseWithClaims(cleanHead, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(api.jwtSecret), nil
	})
	if err != nil {
		respondWithError(w, 401, "Invalid or expired JWT token!")
		log.Printf("Token: %s", err.Error())
		return
	}

	id, _ := token.Claims.GetSubject()
	user, errr := decodeAndValidateUser(r)
	if errr != nil {
		respondWithError(w, 401, errr.Error())
		log.Printf("Decode body: %s", errr.Error())
		return
	}
	// Find and update user
	usr, upd := api.database.UpdateUser(id, user)
	if upd != nil {
		respondWithError(w, 401, upd.Error())
		log.Printf("Database: %s", upd.Error())
		return
	}

	userNoPass := database.UserNoPassword{
		Id: usr.Id,
		Email: usr.Email,
	}

	respondWithJSON(w, 200, userNoPass)
}

func decodeAndValidateUser(r *http.Request) (database.User, error) {
	decoder := json.NewDecoder(r.Body)
	user := database.User{}
	err := decoder.Decode(&user)

	if err != nil {
		return database.User{}, fmt.Errorf("Unable to decode user!")
	}

	msg := validateUser(user)
	if msg != nil {
		return database.User{}, msg
	}

	return user, nil
}

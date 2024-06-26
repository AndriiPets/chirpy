package database

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
)

var ids int = 0

type DB struct {
	path string
	mux *sync.RWMutex
}

type Chirp struct {
	Id int `json:"id"`
	Body string `json:"body"`
}

type User struct {
	Id int `json:"id"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type UserNoPassword struct {
	Id int `json:"id"`
	Email string `json:"email"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users map[int]User `json:"users"`
}

func NewDB(path string) (*DB, error) {
	_, err := os.ReadFile(path)

	if err != nil {
		data, _ := json.Marshal(DBStructure{})
		erro := os.WriteFile(path, data, 0666)
		fmt.Println("hit")
		if erro != nil {
			return &DB{}, erro
		}
	}

	data := DB{
		path: path,
		mux: &sync.RWMutex{},
	}

	return &data, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	ids ++
	id := ids 
	chirp := Chirp{
		Id: id,
		Body: body,
	}
	database, err := db.readDatabase()
	if err != nil {
		return Chirp{}, err
	}
	if len(database.Chirps) == 0 {
		database.Chirps = make(map[int]Chirp)
	}

	database.Chirps[id] = chirp

	db.writeDatabase(database)

	return chirp, nil
	
}

func (db *DB) CreateUser(email, password string) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	ids ++
	id := ids 
	user := User{
		Id: id,
		Email: email,
		Password: password,
	}
	database, err := db.readDatabase()
	if err != nil {
		return User{}, err
	}
	if len(database.Users) == 0 {
		database.Users = make(map[int]User)
	}

	database.Users[id] = user

	db.writeDatabase(database)

	return user, nil
	
}

func (db *DB) GetChirps() ([]Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	chirps := []Chirp{}

	database, err := db.readDatabase()
	if err != nil {
		return chirps, err
	}

	for _, v := range database.Chirps {
		chirps = append(chirps, v)
	}

	return chirps, nil

}

func (db *DB) GetUsers() ([]User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	users := []User{}

	database, err := db.readDatabase()
	if err != nil {
		return users, err
	}

	for _, v := range database.Users {
		users = append(users, v)
	}

	return users, nil

}

func (db *DB) GetChirp(id string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	fmt.Println(id)

	index, err := strconv.Atoi(id)
	if err != nil {
		return Chirp{}, err
	}

	database, err := db.readDatabase()
	if err != nil {
		return Chirp{}, err
	}

	chrp, ok := database.Chirps[index]
	if !ok {
		return Chirp{}, fmt.Errorf("Key does not exist!")
	}

	return chrp, nil

}

func (db *DB) readDatabase() (*DBStructure, error) {
	database := &DBStructure{}
	data, err := os.ReadFile(db.path) 
	if err != nil {
		return database, err
	}

	erro := json.Unmarshal(data, database) 
	if erro != nil {
		return database, erro
	}

	return database, nil
}

func (db *DB)writeDatabase(data *DBStructure) error {
	updataedDB, err := json.Marshal(data) 
	if err != nil {
		return err
	}

	erro := os.WriteFile(db.path, updataedDB, 0666)
	if erro != nil {
		return erro
	}

	return nil
}


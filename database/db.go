package database

import (
	"encoding/json"
	"os"
	"sync"
)

type DB struct {
	path string
	mux *sync.RWMutex
}

type Chirp struct {
	Id int `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

func NewDB(path string) (*DB, error) {
	_, err := os.ReadFile(path)

	if err == os.ErrNotExist {
		data, _ := json.Marshal(DBStructure{})
		erro := os.WriteFile(path, data, 0666)
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

	id := 1
	chirp := Chirp{
		Id: id,
		Body: body,
	}
	database := &DBStructure{}

	data, err := os.ReadFile(db.path) 
		if err != nil {
			return Chirp{}, err
		}
	json.Unmarshal(data, database)
	database.Chirps[id] = chirp

	updataedDB, _ := json.Marshal(database) 
	erro := os.WriteFile(db.path, updataedDB, 0666)
		if erro != nil {
			return Chirp{}, erro
		}

	return chirp, nil
	
}


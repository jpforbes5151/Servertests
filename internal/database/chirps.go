package database

import (
	"os"
)

type Chirp struct {
	Body string `json:"body"`
	ID   int    `json:"id"`
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	// attempt to load overlying DB structure
	dbStructure, err := db.loadDB()
	// if error is thrown, return empty Chirp struct and the error being thrown
	if err != nil {
		return Chirp{}, err
	}

	// make ID for new Chirp after checking current length of Chirp Struct
	id := len(dbStructure.Chirps) + 1
	// making new Chirp Struct from request
	chirp := Chirp{
		Body: body,
		ID:   id,
	}
	dbStructure.Chirps[id] = chirp

	// attempt to write recently made Chirp struct to DB
	err = db.writeDB(dbStructure)
	// return empty Chirp struct and error if an error gets thrown
	if err != nil {
		return Chirp{}, err
	}

	// if no errors get thrown, return the newly made Chirp struct and no error
	return chirp, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	// attempt to load overlying DB structure
	dbStructure, err := db.loadDB()
	// if error is thrown, return empty Chirp struct and the error being thrown
	if err != nil {
		return Chirp{}, err
	}
	// attempt to pull Chirp from DB using ID
	chirp, ok := dbStructure.Chirps[id]

	// if an error gets thrown when trying to retrieve this Chirp Struct, throw an error stating that it does not exist
	if !ok {
		return Chirp{}, os.ErrNotExist
	}

	// if no errors get thrown, return the desired Chirp struct
	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	// // attempt to load overlying DB structure
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	// attempts to make a slice of current Chirp structures and append all entries to this new slice
	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	// return the new slice generated
	return chirps, nil
}

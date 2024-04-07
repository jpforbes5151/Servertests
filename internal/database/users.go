package database

import (
	"errors"
)

type User struct {
	Email          string `json:"email"`
	ID             int    `json:"id"`
	HashedPassword string `json:"hashed_password"`
}

var ErrAlreadyExists = errors.New("already exists")

func (db *DB) CreateUser(email, hashedPassword string) (User, error) {
	if _, err := db.GetUserByEmail(email); !errors.Is(err, ErrNotExist) {
		return User{}, ErrAlreadyExists
	}
	// attempt to load DB
	dbStructure, err := db.loadDB()
	// check if the database has the desired structure
	if err != nil {
		return User{}, err
	}
	// with valid structure, make object of the User structure
	id := len(dbStructure.Users) + 1
	user := User{
		Email:          email,
		ID:             id,
		HashedPassword: hashedPassword,
	}
	dbStructure.Users[id] = user

	// Attempt to write the User struct to DB, if an error gets thrown, return an empty User struct and an error
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	// If no errors get thrown, return the Struct written and no error
	return user, nil
}

func (db *DB) GetUser(id int) (User, error) {
	// try to pull UserID from overlying DB Structure
	dbStructure, err := db.loadDB()
	// if an error gets thrown, return empty user struct and the error being thrown
	if err != nil {
		return User{}, err
	}
	// attempt to retrieve the desired user via ID
	user, ok := dbStructure.Users[id]
	// if error gets thrown, return empty User struct and error being thrown
	if !ok {
		return User{}, ErrNotExist
	}

	// if no errors are thrown, return user that corresponds with ID
	return user, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, ErrNotExist
}

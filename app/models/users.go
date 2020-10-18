package models

import (
	"database/sql"
	"polygnosics/app/services/db"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
)

// User defines the user structures.
type User struct {
	ID       uuid.UUID
	Email    string
	Password string
	Username string
}

// GetUserByEmail returns the user defined by the email and password.
func GetUserByEmail(email string) (User, error) {
	var user User
	queryString := "select BIN_TO_UUID(id), username, email, password from users where email = ?"
	db, err := db.ConnectSystem()
	if err != nil {
		return user, err
	}
	defer db.Close()

	query, err := db.Query(queryString, email)

	if err != nil {
		return user, err
	}
	defer query.Close()

	query.Next()
	query.Scan(&user.ID, &user.Username, &user.Email, &user.Password)

	return user, nil
}

func UserOrEmailExist(email string, username string) (bool, bool, error) {
	var user User
	queryString := "select BIN_TO_UUID(id), username, email, password from users where email = ?"
	db, err := db.ConnectSystem()
	if err != nil {
		return false, false, err
	}
	defer db.Close()

	queryEmail := db.QueryRow(queryString, email)

	err = queryEmail.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	switch {
	case err == sql.ErrNoRows:
		break
	case err != nil:
		return false, false, err
	default:
		return true, false, err
	}

	queryString = "select BIN_TO_UUID(id), username, email, password from users where username = ?"
	queryUser := db.QueryRow(queryString, username)

	err = queryUser.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	switch {
	case err == sql.ErrNoRows:
		break
	case err != nil:
		return false, false, err
	default:
		return false, true, err
	}

	return false, false, nil
}

// CheckPassword compares the password entered by the user with the stored password.
func CheckPassword(email string, password string) (bool, error) {
	user, err := GetUserByEmail(email)
	if err != nil {
		return false, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return false, err
	}
	return true, nil
}

// AddUser adds a user with user name, email, password.
func AddUser(user string, email string, passwd string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwd), 16)
	if err != nil {
		return err
	}

	queryString := "INSERT INTO users (id, username, email, password) VALUES (UUID_TO_BIN(UUID()), ?, ?, ?)"
	db, err := db.ConnectSystem()
	if err != nil {
		return err
	}

	defer db.Close()

	query, err := db.Query(queryString, user, email, hashedPassword)
	if err != nil {
		return err
	}

	defer query.Close()
	return nil
}

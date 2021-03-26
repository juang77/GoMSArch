package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/juang77/GoMSArch/profile-service/app/models"
	"github.com/juang77/GoMSArch/profile-service/config"
)

// OpenConnection method. This method is being used by the main function. For testing the database is being mocked.
func OpenConnection(cnf config.Config) (*sql.DB, error) {
	username := cnf.DBUsername
	password := cnf.DBPassword
	host := cnf.DBHost
	port := cnf.DBPort
	database := cnf.Database

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true", username, password, host, port, database)

	log.Debugf("Connect to : %v", dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, ErrCanNotConnectWithDatabase
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		return nil, ErrCanNotConnectWithDatabase
	}
	return db, nil
}

// CloseConnection method
func CloseConnection(db *sql.DB) {
	db.Close()
}

// GetUserByID returns an models.User identified by it's ID or a ErrUserNotFound error when the user cannot be found.
func GetUserByID(db *sql.DB, ID int64) (models.UserResponse, error) {
	rows, err := db.Query("CALL usp_get_user_by_id(?)", ID)
	if err != nil {
		return models.UserResponse{}, err
	}

	if rows.Next() {
		var user_id int64
		var user_name string
		var createdAt time.Time
		var user_mail string
		var user_mobile string
		err = rows.Scan(&user_id, &user_name, &createdAt, &user_mail, &user_mobile)
		if err != nil {
			return models.UserResponse{}, err
		}

		return models.UserResponse{ID: user_id, Username: user_name, CreatedAt: createdAt, Email: user_mail, Mobile: user_mobile}, nil
	}
	return models.UserResponse{}, ErrUserNotFound
}

// CreateUser create an user in the database and returns the ID of the user being inserted. This method returns a ErrUsernameIsNotUnique or ErrEmailIsNotUnique when the username or email of an user is not unique.
func CreateUser(db *sql.DB, user *models.UserCreate) (int64, error) {
	// Insert
	res, err := db.Query("CALL usp_create_user(?, ?, ?, ?)", user.Username, user.Email, user.Mobile, user.Hash)
	if err != nil {
		log.Errorf("Error inserting")
		log.Error(err)
		return 0, err
	}

	response := make([]models.SPResponse, 0)

	if res.Next() {
		var spResponse models.SPResponse
		/* Leemos el registro */
		res.Scan(&spResponse.Msg, &spResponse.Id)
		response = append(response, spResponse)
	}

	if response[0].Msg == "Duplicated Name" {
		return 0, ErrUsernameIsNotUnique
	}

	if response[0].Msg == "Duplicated Email" {
		return 0, ErrEmailIsNotUnique
	}

	return int64(response[0].Id), nil
}

// UpdateUser updates the username and email of an user. (note: this method does not check if user is authorized to update this row)
func UpdateUser(db *sql.DB, user *models.UserResponse) (int64, error) {
	//Update
	res, err := db.Query("CALL usp_update_user(?, ?. ?)", user.ID, user.Email, user.Mobile)
	if err != nil {
		log.Errorf("Error Updating")
		log.Error(err)
		return 0, err
	}

	response := make([]models.SPResponse, 0)

	if res.Next() {
		var spResponse models.SPResponse
		/* Leemos el registro */
		res.Scan(&spResponse.Msg, &spResponse.Id)
		response = append(response, spResponse)
	}

	if response[0].Msg == "Duplicated Email" {
		return 0, ErrEmailIsNotUnique
	}

	return user.ID, nil
}

// DeleteUser deletes an user from the database. Method does not check if the caller is authorized to perform this action. Method returns the number of rows affected by query. (should be 1)
func DeleteUser(db *sql.DB, user *models.UserResponse) (int64, error) {
	if user.ID > 0 {
		res, err := db.Query("CALL usp_delete_user(?)", user.ID)
		if err != nil {
			log.Errorf("Error inserting")
			log.Error(err)
			return 0, err
		}
		response := make([]models.SPResponse, 0)

		if res.Next() {
			var spResponse models.SPResponse
			/* Leemos el registro */
			res.Scan(&spResponse.Msg, &spResponse.Id)
			response = append(response, spResponse)
		}

		if response[0].Msg == "Not Exist" {
			return 0, errors.New("user id is empty")
		}
	}
	return user.ID, nil
}

// GetUsers returns a list of all database-users. Note: Consider implementing a paging function because this method returns EVERY users at once.
func GetUsers(db *sql.DB) ([]models.UserResponse, error) {

	rows, err := db.Query("Call usp_get_users()")
	if err != nil {
		return nil, err
	}
	persons := make([]models.UserResponse, 0)

	for rows.Next() {
		var user_id int64
		var user_name string
		var user_mail string
		var user_mobile string
		var createdAt time.Time
		err = rows.Scan(&user_id, &user_name, &user_mail, &user_mobile, &createdAt)
		if err != nil {
			return nil, err
		}
		persons = append(persons, models.UserResponse{ID: user_id, Username: user_name, Mobile: user_mail, CreatedAt: createdAt})
	}
	return persons, nil
}

// ErrEmailIsNotUnique error is the email is not unique
var ErrEmailIsNotUnique = errors.New("email must be unique")

// ErrUsernameIsNotUnique error if the username is not unique
var ErrUsernameIsNotUnique = errors.New("username must be unique")

// ErrUserNotFound error if user does not exist in database
var ErrUserNotFound = errors.New("user does not exist")

// ErrCanNotConnectWithDatabase error if database is unreachable
var ErrCanNotConnectWithDatabase = errors.New("can not connect with database")

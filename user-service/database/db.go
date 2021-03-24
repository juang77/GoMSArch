package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/juang77/GoMSArch/user-service/app/models"
	"github.com/juang77/GoMSArch/user-service/config"
	log "github.com/sirupsen/logrus"
)

// OpenConnection opens the connection to the database
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

	return db, nil
}

// CloseConnection closes the connection to the database
func CloseConnection(db *sql.DB) {
	db.Close()
}

// GetUsers return the models.User object based on the username
func GetUsers(db *sql.DB) ([]models.User, error) {
	// Query the database
	rows, err := db.Query("CALL usp_get_all_users()")
	if err != nil {
		return []models.User{}, err
	}

	resulted := []models.User{}

	/* Bucle para recorrer todos los registros */
	for rows.Next() {
		var user models.User
		/* Leemos el registro */
		rows.Scan(&user.ID, &user.Name, &user.Email, &user.Mobile)

		resulted = append(resulted, user)
	}
	if len(resulted) <= 0 {
		return resulted, ErrUserNotFound
	} else {
		return resulted, nil
	}

}

// ErrUserNotFound error if user does not exist in database
var ErrUserNotFound = errors.New("does not exist any user")

// ErrCanNotConnectWithDatabase error if database is unreachable
var ErrCanNotConnectWithDatabase = errors.New("can not connect with database")

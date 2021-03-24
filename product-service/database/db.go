package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/juang77/GoMSArch/product-service/app/models"
	"github.com/juang77/GoMSArch/product-service/config"
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

// GetProducts return the models.Product object based
func GetProducts(db *sql.DB) ([]models.Product, error) {
	// Query the database
	rows, err := db.Query("CALL usp_get_all_products()")
	if err != nil {
		return []models.Product{}, err
	}

	resulted := []models.Product{}

	/* Bucle para recorrer todos los registros */
	for rows.Next() {
		var product models.Product
		/* Leemos el registro */
		rows.Scan(&product.ID, &product.Name, &product.Description, &product.Prise)

		resulted = append(resulted, product)
	}
	if len(resulted) <= 0 {
		return resulted, ErrProductNotFound
	} else {
		return resulted, nil
	}

}

// ErrProductNotFound error if product does not exist in database
var ErrProductNotFound = errors.New("does not exist any user")

// ErrCanNotConnectWithDatabase error if database is unreachable
var ErrCanNotConnectWithDatabase = errors.New("can not connect with database")

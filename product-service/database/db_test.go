package db

import (
	"testing"

	"github.com/juang77/GoMSArch/product-service/app/models"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Define user
	product := getTestProduct()
	product1 := getTestProduct1()
	product2 := getTestProduct2()
	// Expected rows
	rows := sqlmock.NewRows([]string{"ID", "Name", "Description", "Prise"}).AddRow(product.ID, product.Name, product.Description, product.Price).AddRow(product1.ID, product1.Name, product1.Description1, product1.Price).AddRow(product2.ID, product2.Name, product2.Description1, product2.Price)

	// define expectations
	// Expectation: check for unique username
	mock.ExpectQuery("CALL usp_get_all_users()").WillReturnRows(rows)

	// Execute the method
	if _, err := GetProducts(db); err != nil {
		t.Errorf("there was an unexpected error: %s", err)
	}

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func getTestProduct() *models.Product {
	product := &models.Product{}
	product.ID = 1
	product.Name = "Comida De Perro 10"
	product.Description = "Comida de perro test 10 para hacer la prueba"
	product.Price = 10000
	return product
}

func getTestProduct1() *models.Product {
	product := &models.Product{}
	product.ID = 2
	product.Name = "Comida De Perro 20"
	product.Description = "Comida de perro test 20 para hacer la prueba"
	product.Price = 20000
	return product
}

func getTestProduct2() *models.Product {
	product := &models.Product{}
	product.ID = 3
	product.Name = "Comida De Perro 30"
	product.Description = "Comida de perro test 30 para hacer la prueba"
	product.Price = 30000
	return product
}

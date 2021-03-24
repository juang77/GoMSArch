package db

import (
	"testing"

	"github.com/juang77/GoMSArch/user-service/app/models"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Define user
	user := getTestUser()
	user1 := getTestUser1()
	// Expected rows
	rows := sqlmock.NewRows([]string{"ID", "Name", "Email", "Mobile"}).AddRow(user.ID, user.Name, user.Email, user.Mobile).AddRow(user1.ID, user1.Name, user1.Email, user1.Mobile)

	// define expectations
	// Expectation: check for unique username
	mock.ExpectQuery("CALL usp_get_all_users()").WillReturnRows(rows)

	// Execute the method
	if _, err := GetUsers(db); err != nil {
		t.Errorf("there was an unexpected error: %s", err)
	}

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func getTestUser() *models.User {
	user := &models.User{}
	user.ID = 1
	user.Name = "JUAN GABRIEL GOMEZ BARON"
	user.Email = "juang77@hotmail.com"
	user.Mobile = "0313166937870"
	return user
}

func getTestUser1() *models.User {
	user := &models.User{}
	user.ID = 2
	user.Name = "ANGELA PATRICIA PARRA LOPEZ"
	user.Email = "angela_parra76@hotmail.com"
	user.Mobile = "0313212358748"
	return user
}

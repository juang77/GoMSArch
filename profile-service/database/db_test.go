package db

import (
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/juang77/GoMSArch/profile-service/app/models"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetUserByID(t *testing.T) {
	user := getTestUser()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Expectation: select user where id == user.ID
	timeNow := time.Now().UTC()
	selectByIDRows := sqlmock.NewRows([]string{"user_id", "user_name", "createdAt", "user_mail", "user_mobile"}).AddRow(user.ID, user.Username, timeNow, user.Email, user.Mobile)
	mock.ExpectQuery("CALL usp_get_user_by_id(?)").WithArgs(user.ID).WillReturnRows(selectByIDRows)

	// Execute the method
	if _, err := GetUserByID(db, user.ID); err != nil {
		t.Errorf("there was an unexpected error: %s", err)
	}

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateUser(t *testing.T) {
	// Define user
	user := getTestUserForCreation()
	//Define Ok Response
	response := getOKresponse()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Expected rows
	rows := sqlmock.NewRows([]string{"msg", "id"}).AddRow(response.Msg, response.Id)

	// Expectation: insert into database
	mock.ExpectQuery("CALL usp_create_user(?)").WithArgs(user.Username, user.Email, user.Mobile, user.Hash).WillReturnRows(rows)

	// Execute the method
	if _, err := CreateUser(db, user); err != nil {
		t.Errorf("there was an unexpected error: %s", err)
	}

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestUpdateUser
func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Define user
	user := getTestUser()

	// Expected rows
	response := getOKresponse()
	rows := sqlmock.NewRows([]string{"msg", "id"}).AddRow(response.Msg, response.Id)

	// Expectation: insert into database
	mock.ExpectQuery("CALL usp_update_user(?)").WithArgs(user.ID, user.Email, user.Mobile).WillReturnRows(rows)

	// Execute the method
	if _, err := UpdateUser(db, user); err != nil {
		t.Errorf("there was an unexpected error: %s", err)
	}

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestDeleteUser
func TestDeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Define user
	user := getTestUser()

	// Expected rows
	response := getOKresponse()
	rows := sqlmock.NewRows([]string{"msg", "id"}).AddRow(response.Msg, response.Id)

	// Expectation: insert into database
	mock.ExpectQuery("CALL usp_delete_user(?)").WithArgs(user.ID).WillReturnRows(rows)

	// Execute the method
	if _, err := DeleteUser(db, user); err != nil {
		t.Errorf("there was an unexpected error: %s", err)
	}

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestGetUsers
func TestGetUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	user1 := getTestUser()
	user2 := getTestUser()
	user2.ID = 2

	// Expectation: insert into database
	timeNow := time.Now().UTC()
	selectRows := sqlmock.NewRows([]string{"user_id", "user_name", "user_mail", "user_mobile", "createdAt"}).AddRow(user1.ID, user1.Username, user1.Email, user1.Mobile, timeNow).AddRow(user2.ID, user2.Username, user2.Email, user2.Mobile, timeNow)
	mock.ExpectQuery("Call usp_get_users()").WillReturnRows(selectRows)

	// Execute the method
	if _, err := GetUsers(db); err != nil {
		t.Errorf("there was an unexpected error: %s", err)
	}

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func getTestUserForCreation() *models.UserCreate {
	user := &models.UserCreate{}
	user.ID = 4
	user.Email = "nicolasgp2018@hotmail.com"
	user.Password = "12345678"
	user.Username = "NICOLAS GOMEZ PARRA"
	user.CreatedAt = time.Now()
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.Hash = string(hash)
	return user
}

func getTestUser() *models.UserResponse {
	user := &models.UserResponse{}
	user.ID = 1
	user.Email = "juang77@hotmail.com"
	user.Username = "JUAN GABRIEL GOMEZ BARON"
	user.CreatedAt = time.Now()
	user.Mobile = "0313166937870"
	return user
}

func getOKresponse() *models.SPResponse {
	spResponse := &models.SPResponse{}
	spResponse.Msg = "OK"
	spResponse.Id = 1
	return spResponse
}

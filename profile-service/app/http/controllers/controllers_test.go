package controllers

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/juang77/GoMSArch/profile-service/app/models"
	"github.com/juang77/GoMSArch/profile-service/config"
	"github.com/urfave/negroni"
	"golang.org/x/crypto/bcrypt"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type TestHash struct{}

func (a TestHash) Match(v driver.Value) bool {
	_, ok := v.(string)
	return ok
}

// // Test creating an user.
// func TestCreateUser(t *testing.T) {
// 	cnf := config.Config{}
// 	cnf.SecretKey = "Juancho"

// 	user := getTestUserForCreation()

// 	json, _ := json.Marshal(user)

// 	req, err := http.NewRequest("POST", "http://localhost/users", bytes.NewBuffer(json))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	res := httptest.NewRecorder()

// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()

// 	//Define Ok Response
// 	responseOK := getOKresponse()

// 	// Expected rows
// 	rows := sqlmock.NewRows([]string{"msg", "id"}).AddRow(responseOK.Msg, responseOK.Id)

// 	// Expectation: insert into database
// 	mock.ExpectQuery("CALL usp_create_user(?)").WithArgs(user.Username, user.Email, user.Mobile, user.Hash).WillReturnRows(rows)

// 	// Expectation: select user where id == user.ID
// 	timeNow := time.Now().UTC()
// 	selectByIDRows := sqlmock.NewRows([]string{"user_id", "user_name", "createdAt", "user_mail", "user_mobile"}).AddRow(user.ID, user.Username, timeNow, user.Email, user.Mobile)
// 	mock.ExpectQuery("CALL usp_get_user_by_id(?)").WithArgs(user.ID).WillReturnRows(selectByIDRows)

// 	handler := CreateUserHandler(db, cnf)
// 	handler(res, req, nil)

// 	// Make sure expectations are met
// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("there were unfulfilled expectations: %s", err)
// 		fmt.Println(err)
// 	}

// 	// Make sure response is alright
// 	type Token struct {
// 		Token     string              `json:"token"`
// 		ExpiresOn string              `json:"expires_on"`
// 		User      models.UserResponse `json:"user"`
// 	}

// 	response := &Token{}
// 	err = decodeJSON(res.Body, response)
// 	if err != nil {
// 		t.Fatal(errors.New("Bad json"))
// 	}
// 	if response.User.ID < 1 {
// 		t.Errorf("Expected user ID greater than 0 but got %v", response.User.ID)
// 	}

// 	if response.User.Username != user.Username {
// 		t.Errorf("Expected username to be %v but got %v", user.Username, response.User.Username)
// 	}
// 	if response.User.Email != user.Email {
// 		t.Errorf("Expected username to be %v but got %v", user.Email, response.User.Email)
// 	}

// 	if res.Result().StatusCode != 200 {
// 		t.Errorf("Expected statuscode to be 200 but got %v", res.Result().StatusCode)
// 	}
// }

// Test creating an user when a bad json string is provided. We expect an error message.
func TestBadJson(t *testing.T) {
	cnf := config.Config{}
	cnf.SecretKey = "Juancho"

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	req, err := http.NewRequest("POST", "http://localhost/users", bytes.NewBuffer([]byte("{")))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handler := CreateUserHandler(db, cnf)
	handler(res, req, nil)

	actual := res.Body.String()
	expected := "{\"message\":\"Bad json\"}"
	if expected != actual {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
}

// Test creating a user without providing a username. We expect an error message.
func TestCreateUserWithoutUsername(t *testing.T) {
	cnf := config.Config{}
	cnf.SecretKey = "ABCDEF"

	user := &models.UserCreate{}

	json, _ := json.Marshal(user)

	req, err := http.NewRequest("POST", "http://localhost/users", bytes.NewBuffer(json))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	handler := CreateUserHandler(db, cnf)
	handler(res, req, nil)

	actual := res.Body.String()
	expected := "{\"message\":\"username is too short\"}"
	if expected != actual {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
}

// Test creating an user without providing a password. We expect an error message.
func TestCreateUserWithoutPassword(t *testing.T) {
	cnf := config.Config{}
	cnf.SecretKey = "ABCDEF"

	user := &models.UserCreate{}
	user.Email = "test@example.com"
	user.Username = "username"

	json, _ := json.Marshal(user)

	req, err := http.NewRequest("POST", "http://localhost/users", bytes.NewBuffer(json))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	handler := CreateUserHandler(db, cnf)
	handler(res, req, nil)

	actual := res.Body.String()
	expected := "{\"message\":\"password is to short\"}"
	if expected != actual {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
}

// Test creating an user without providing an email address. We expect an error message.
func TestCreateUserWithoutEmail(t *testing.T) {
	cnf := config.Config{}
	cnf.SecretKey = "ABCDEF"

	user := &models.UserCreate{}
	user.Username = "NICOLAS GOMEZ PARRA"

	json, _ := json.Marshal(user)

	req, err := http.NewRequest("POST", "http://localhost/users", bytes.NewBuffer(json))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	handler := CreateUserHandler(db, cnf)
	handler(res, req, nil)

	actual := res.Body.String()
	expected := "{\"message\":\"email address is to short\"}"
	if expected != actual {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
}

// Helper function to decode a json string to an interface
func decodeJSON(r io.Reader, target interface{}) error {
	err := json.NewDecoder(r).Decode(target)
	if err != nil {
		fmt.Printf("json decoder error occurred: %v \n", err.Error())
		return err
	}
	return nil
}

// Test deleting an user.
func TestDeleteUser(t *testing.T) {
	cnf := config.Config{}
	cnf.SecretKey = "Juancho"

	user := &models.UserCreate{}
	user.ID = 1
	user.Email = "username@example.com"
	user.Password = "password"
	user.Username = "username"

	// Create JWT object with claims
	expiration := time.Now().Add(time.Hour * 24 * 31).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"iat": time.Now().Unix(),
		"exp": expiration,
	})

	// Generate a signed token
	secretKey := cnf.SecretKey
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		t.Error(err)
		return
	}
	json, _ := json.Marshal(user)
	req, err := http.NewRequest("DELETE", "http://localhost/users?token="+tokenString, bytes.NewBuffer(json))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Expected rows
	response := getOKresponse()
	rows := sqlmock.NewRows([]string{"msg", "id"}).AddRow(response.Msg, response.Id)

	// Expectation: insert into database
	mock.ExpectQuery("CALL usp_delete_user(?)").WithArgs(user.ID).WillReturnRows(rows)

	handler := DeleteUserHandler(db, cnf)
	handler(res, req, nil)

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	if res.Result().StatusCode != 200 {
		t.Errorf("Expected statuscode to be 200 but got %v", res.Result().StatusCode)
	}
}

// Test updating an user.
func TestUpdateUser(t *testing.T) {
	cnf := config.Config{}
	cnf.SecretKey = "Juancho"

	user := getTestUser()
	tokenString := getTokenString(cnf, user, t)

	json, _ := json.Marshal(user)
	req, err := http.NewRequest("PUT", "http://localhost/users?token="+tokenString, bytes.NewBuffer(json))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Expected rows
	response := getOKresponse()
	rows := sqlmock.NewRows([]string{"msg", "id"}).AddRow(response.Msg, response.Id)
	mock.ExpectQuery("CALL usp_update_user(?)").WithArgs(user.ID, user.Email, user.Mobile).WillReturnRows(rows)

	handler := UpdateUserHandler(db, cnf)
	handler(res, req, nil)

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	if res.Result().StatusCode != 200 {
		t.Errorf("Expected statuscode to be 200 but got %v", res.Result().StatusCode)
	}
}

// Get an user by it's index.
func TestGetUserByIndex(t *testing.T) {
	// Mock user object
	user := getTestUser()

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Expectation: select user where id == user.ID
	timeNow := time.Now().UTC()
	selectByIDRows := sqlmock.NewRows([]string{"user_id", "user_name", "createdAt", "user_mail", "user_mobile"}).AddRow(user.ID, user.Username, timeNow, user.Email, user.Mobile)
	mock.ExpectQuery("CALL usp_get_user_by_id(?)").WithArgs(user.ID).WillReturnRows(selectByIDRows)

	// Router
	r := mux.NewRouter()
	r.Handle("/user/{id}", negroni.New(
		negroni.HandlerFunc(UserByIndexHandler(db)),
	))

	// Server
	ts := httptest.NewServer(r)
	defer ts.Close()

	// Do Request
	url := ts.URL + "/user/" + fmt.Sprintf("%v", user.ID)
	res, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Make sure response statuscode expectation is met
	if res.StatusCode != 200 {
		t.Errorf("Expected statuscode to be 200 but got %v", res.StatusCode)
	}
}

// In this test we'll login as user1 and we try to change user2. This is not allowed therefore we expect an error.
func TestTryUpdateOtherUser(t *testing.T) {
	cnf := config.Config{}
	cnf.SecretKey = "Juancho"

	user1 := getTestUser()
	user1.ID = 1
	user1.Username = "user1"

	user2 := getTestUser()
	user2.ID = 2
	user2.Username = "user2"

	tokenString := getTokenString(cnf, user1, t)

	json, _ := json.Marshal(user2)
	req, err := http.NewRequest(http.MethodPut, "http://localhost/users?token="+tokenString, bytes.NewBuffer(json))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handler := UpdateUserHandler(nil, cnf)
	handler(res, req, nil)

	// Make sure expectations are met
	expected := `{"message":"you can only change your own user object"}`
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			res.Body.String(), expected)
	}
	if res.Result().StatusCode != 400 {
		t.Errorf("Expected statuscode to be 400 but got %v", res.Result().StatusCode)
	}
}

// In this test we'll login as user1 and we try to delete user2. This is not allowed therefore we expect an error.
func TestTryDeleteOtherUser(t *testing.T) {
	cnf := config.Config{}
	cnf.SecretKey = "Juancho"

	user1 := getTestUser()
	user1.ID = 1
	user1.Username = "user1"

	user2 := getTestUser()
	user2.ID = 2
	user2.Username = "user2"

	tokenString := getTokenString(cnf, user1, t)

	json, _ := json.Marshal(user2)
	req, err := http.NewRequest(http.MethodDelete, "http://localhost/users?token="+tokenString, bytes.NewBuffer(json))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handler := UpdateUserHandler(nil, cnf)
	handler(res, req, nil)

	// Make sure expectations are met
	expected := `{"message":"you can only change your own user object"}`
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			res.Body.String(), expected)
	}
	if res.Result().StatusCode != 400 {
		t.Errorf("Expected statuscode to be 400 but got %v", res.Result().StatusCode)
	}
}

func getTestUser() *models.UserCreate {
	user := &models.UserCreate{}
	user.ID = 4
	user.Email = "nicolasgp2018@hotmail.com"
	user.Username = "NICOLAS GOMEZ PARRA"
	user.CreatedAt = time.Now()
	user.Password = "12345678"
	user.Mobile = "0313166937871"
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.Hash = string(hash)
	return user
}

func getTestUserForCreation() *models.UserCreate {
	user := &models.UserCreate{}
	user.ID = 4
	user.Email = "nicolasgp2018@hotmail.com"
	user.Password = "12345678"
	user.Username = "NICOLAS GOMEZ PARRA"
	user.CreatedAt = time.Now()
	user.Mobile = "0313166937871"
	user.Hash = ""
	return user
}

func getOKresponse() *models.SPResponse {
	spResponse := &models.SPResponse{}
	spResponse.Msg = "OK"
	spResponse.Id = 1
	return spResponse
}

func getTokenString(cnf config.Config, user *models.UserCreate, t *testing.T) string {
	expiration := time.Now().Add(time.Hour * 24 * 31).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"iat": time.Now().Unix(),
		"exp": expiration,
	})

	// Generate a signed token
	secretKey := cnf.SecretKey
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		t.Error(err)
		return ""
	}
	return tokenString
}

package routes

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/juang77/GoMSArch/profile-service/app/models"
	"github.com/juang77/GoMSArch/profile-service/config"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type TestHash struct{}

func (a TestHash) Match(v driver.Value) bool {
	_, ok := v.(string)
	return ok
}

func TestOPTIONSUsers(t *testing.T) {
	// Router
	r := InitRoutes(nil, config.Config{})
	res := httptest.NewRecorder()

	// Do Request
	req, err := http.NewRequest(http.MethodOptions, "/users", nil)
	if err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(res, req)

	// Make sure response statuscode expectation is met
	if res.Result().StatusCode != 200 {
		t.Errorf("Expected statuscode to be 200 but got %v", res.Result().StatusCode)
	} else {
		t.Logf("Result statuscode %v. (As expected)", res.Result().StatusCode)
	}
}

func TestPUTUsers(t *testing.T) {

	user := getTestUser()

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Expected rows
	response := getOKresponse()
	rows := sqlmock.NewRows([]string{"msg", "id"}).AddRow(response.Msg, response.Id)

	// Expectation: insert into database
	mock.ExpectQuery("CALL usp_update_user(?)").WithArgs(user.ID, user.Email, user.Mobile).WillReturnRows(rows)

	// Mock config
	cnf := config.Config{}
	cnf.SecretKey = "Juancho"

	// Get token string
	tokenString := getTokenString(cnf, user, t)

	json, _ := json.Marshal(user)
	res := doRequest(db, cnf, http.MethodPut, "/users?token="+tokenString, bytes.NewBuffer(json), t)

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Make sure response statuscode expectation is met
	if res.Result().StatusCode != 200 {
		t.Errorf("Expected statuscode to be 200 but got %v", res.Result().StatusCode)
	} else {
		t.Logf("Result statuscode %v. (As expected)", res.Result().StatusCode)
	}
}

func TestDELETEUsers(t *testing.T) {
	user := getTestUser()

	// Mock database
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

	// Mock config
	cnf := config.Config{}
	cnf.SecretKey = "Juancho"

	// Get token string
	tokenString := getTokenString(cnf, user, t)

	json, _ := json.Marshal(user)
	res := doRequest(db, cnf, http.MethodDelete, "/users?token="+tokenString, bytes.NewBuffer(json), t)

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Make sure response statuscode expectation is met
	if res.Result().StatusCode != 200 {
		t.Errorf("Expected statuscode to be 200 but got %v", res.Result().StatusCode)
	} else {
		t.Logf("Result statuscode %v. (As expected)", res.Result().StatusCode)
	}
}

// func TestPOSTUsers(t *testing.T) {
// 	user := getTestUser()

// 	// Mock database
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()

// 	// Expected rows
// 	rows := sqlmock.NewRows([]string{"count (*)"})

// 	// Expectation: check for unique username
// 	mock.ExpectQuery("SELECT (.+) FROM users WHERE").WithArgs(user.Username).WillReturnRows(rows)

// 	// Expectation: check for unique email
// 	mock.ExpectQuery("SELECT (.+) FROM users WHERE").WithArgs(user.Email).WillReturnRows(rows)

// 	// Expectation: insert into database
// 	mock.ExpectExec("INSERT INTO users").WithArgs(user.Username, user.Email, TestHash{}).WillReturnResult(sqlmock.NewResult(1, 1))

// 	timeNow := time.Now().UTC()
// 	selectByIDRows := sqlmock.NewRows([]string{"id", "username", "createdAt", "email"}).AddRow(user.ID, user.Username, timeNow, user.Email)
// 	mock.ExpectQuery("SELECT (.+) FROM users WHERE").WithArgs(user.ID).WillReturnRows(selectByIDRows)

// 	// Mock config
// 	cnf := config.Config{}
// 	cnf.SecretKey = "ABC"

// 	// Get token string
// 	tokenString := getTokenString(cnf, user, t)

// 	json, _ := json.Marshal(user)
// 	res := doRequest(db, cnf, http.MethodPost, "/users?token="+tokenString, bytes.NewBuffer(json), t)

// 	// Make sure expectations are met
// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("there were unfulfilled expectations: %s", err)
// 	}

// 	// Make sure response statuscode expectation is met
// 	if res.Result().StatusCode != 200 {
// 		t.Logf(string(res.Body.Bytes()))
// 		t.Errorf("Expected statuscode to be 200 but got %v", res.Result().StatusCode)
// 	} else {
// 		t.Logf("Result statuscode %v. (As expected)", res.Result().StatusCode)
// 	}
// }

func TestGETUserByID(t *testing.T) {
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

	// Mock config
	cnf := config.Config{}
	cnf.SecretKey = "Juancho"

	// Get token string
	tokenString := getTokenString(cnf, user, t)

	json, _ := json.Marshal(user)
	url := "/user/" + fmt.Sprintf("%v", user.ID) + "?token=" + tokenString
	res := doRequest(db, cnf, http.MethodGet, url, bytes.NewBuffer(json), t)

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Make sure response statuscode expectation is met
	if res.Result().StatusCode != 200 {
		t.Errorf("Expected statuscode to be 200 but got %v", res.Result().StatusCode)
	} else {
		t.Logf("Result statuscode %v. (As expected)", res.Result().StatusCode)
	}
}

func doRequest(db *sql.DB, cnf config.Config, method string, url string, body *bytes.Buffer, t *testing.T) *httptest.ResponseRecorder {
	r := InitRoutes(db, cnf)
	res := httptest.NewRecorder()
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(res, req)
	return res
}

func getTestUser() *models.UserCreate {
	user := &models.UserCreate{}
	user.ID = 1
	user.Email = "juang77@hotmail.com"
	user.Username = "JUAN GABRIEL GOMEZ BARON"
	user.CreatedAt = time.Now()
	user.Mobile = "0313166937870"
	user.Hash = "TempFakeHash"
	return user
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

func getOKresponse() *models.SPResponse {
	spResponse := &models.SPResponse{}
	spResponse.Msg = "OK"
	spResponse.Id = 1
	return spResponse
}

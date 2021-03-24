package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/juang77/GoMSArch/user-service/config"
)

func TestGetUsers(t *testing.T) {
	// MOCK SERVER
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL+"/users", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Add Rows
	SelectRows := sqlmock.NewRows([]string{"ID", "Name", "Email", "Mobile"}).AddRow(1, "JUAN GABRIEL GOMEZ BARON", "juang77@hotmail.com", "0313166937870").AddRow(2, "ANGELA PATRICIA PARRA LOPEZ", "angela_parra76@hotmail.com", "0313212358748")
	mock.ExpectQuery("CALL usp_get_all_users()").WillReturnRows(SelectRows)

	// Mock config
	cnf := config.Config{}
	cnf.UsersServiceBaseUrl = ts.URL + "/"
	handler := GetUsersHandler(db, cnf)
	handler(res, req, nil)

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	if res.Result().StatusCode != 200 {
		t.Errorf("Expected statuscode to be 200 but got %v", res.Result().StatusCode)

		t.Errorf(res.Body.String())
	}
}

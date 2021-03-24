package routes

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juang77/GoMSArch/product-service/config"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetProducts(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer ts.Close()

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Add Rows
	SelectRows := sqlmock.NewRows([]string{"ID", "Name", "Description", "Prise"}).AddRow(1, "Comida De Perro 10", "Comida de perro test 10 para hacer la prueba", 10000).AddRow(2, "Comida De Perro 20", "Comida de perro test 20 para hacer la prueba", 20000).AddRow(3, "Comida De Perro 30", "Comida de perro test 30 para hacer la prueba", 30000)
	mock.ExpectQuery("CALL usp_get_all_products()").WillReturnRows(SelectRows)

	cnf := config.Config{}
	cnf.UsersServiceBaseUrl = ts.URL + "/"

	res := doRequest(db, cnf, "GET", ts.URL+"/products", bytes.NewBuffer([]byte("")), t)

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Make sure response statuscode expectation is met
	if res.Result().StatusCode != 200 {
		t.Errorf("Expected statuscode to be 200 but got %v", res.Result().StatusCode)
		t.Errorf(res.Body.String())
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

package main

import (
	"bytes"
	"github.com/antoniocs/go_api/pkg/Config"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//const DbCreation = "DROP DATABASE IF EXISTS `go_api_testing`; CREATE DATABASE IF NOT EXISTS `go_api_testing` USE `go_api_testing`;"

//func createTestDatabase(conn *sql.DB) {
//	conn.Exec(DbCreation)
//
//	//go:embed hello.txt
//	var s string
//
//
//}

func cleanDatabase() {

}


func TestApp(t *testing.T) {
	//@TODO Possibly pass other setting for a custom db
	settings := Config.NewSettings()
	router := InitApp(settings)
	//@TODO Move this to another function
	conn := createDbConnection(settings)

	_, err := conn.Exec("TRUNCATE TABLE guest")
	if err != nil {
		panic("Unable to clear guest table" + err.Error())
	}

	type testData struct {
		Endpoint string
		Method string
		Data []byte
		ExpectedStatusCode int
		ExpectedBodyResponse string
	}

	//NOTE: The order does matter
	//If more time was available I would implement more isolated tests
	var calls = []testData{
		{
			"/api/guestbook/testUser",
			"POST",
			[]byte(`{"table": 1,"accompanying_guests": 2}`),
			http.StatusCreated,
			`{"name":"testUser"}`,
		},
		{
			"/api/guestbook/list",
			"GET",
			[]byte{},
			http.StatusOK,
			`[{"name":"testUser","accompanyingGuests":2,"tableId":1}]`,
		},
		{
			"/api/guestbook/testUser",
			"PUT",
			[]byte(`{"accompanying_guests": 3}`),
			http.StatusAccepted,
			`{"name":"testUser"}`,
		},
		{
			"/api/guestbook/emptySeats",
			"GET",
			[]byte{},
			http.StatusOK,
			`{"emptySeats":18}`,
		},
		{
			"/api/guestbook/testUser",
			"DELETE",
			[]byte{},
			http.StatusAccepted,
			`{"name":"testUser"}`,
		},
	}

	for _, td := range calls {
		req, err  := http.NewRequest(td.Method, td.Endpoint, bytes.NewBuffer(td.Data))

		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != td.ExpectedStatusCode {
			t.Errorf("handler returned wrong status code: got %v want %v", status, td.ExpectedStatusCode)
		}

		if len(td.ExpectedBodyResponse) > 0 {
			if strings.TrimSuffix(rr.Body.String(), "\n") != td.ExpectedBodyResponse {
				t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), td.ExpectedBodyResponse)
			}
		}
	}
}

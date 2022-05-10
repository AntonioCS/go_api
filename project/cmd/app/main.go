package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/antoniocs/go_api/pkg/Config"
	"github.com/antoniocs/go_api/pkg/GuestBook"
	"github.com/antoniocs/go_api/pkg/GuestBook/DataManager"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)


type handleRequests struct {
	guestBook GuestBook.GuestBook
	Router *mux.Router
}

func NewHandleRequests(guestBook GuestBook.GuestBook, port string) *mux.Router {
	hr := handleRequests{
		guestBook,
		mux.NewRouter(),
	}

	hr.setupEndpoints()

	return hr.Router
}

func (hr handleRequests) setupEndpoints() {
	hr.Router.HandleFunc("/api/guestbook/list", hr.handleGuestList)
	hr.Router.HandleFunc("/api/guestbook/{name}", hr.handleGuestAdd).Methods("POST")
	hr.Router.HandleFunc("/api/guestbook/{name}", hr.handleGuestArrives).Methods("PUT")
	hr.Router.HandleFunc("/api/guestbook/{name}", hr.handleGuestRemove).Methods("DELETE")
	hr.Router.HandleFunc("/api/guestbook/emptySeats", hr.handleEmptySeats)
}

func (hr handleRequests) handleGuestList(w http.ResponseWriter, r *http.Request) {
	hr.guestBook.List()

	result := hr.guestBook.List()

	outputData(result, w, http.StatusOK)
}

func (hr handleRequests) handleGuestAdd(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	bodyData := struct {
		Table    int `json:"table"`
		AccompanyingGuests   int `json:"accompanying_guests"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&bodyData)
	if err != nil {
		errorResponse(w, err.Error())
		return
	}

	nameResponse(hr.guestBook.Add(name, bodyData.AccompanyingGuests, bodyData.Table), name, w, http.StatusCreated, "Unable to add guest: " + name)
}

func (hr handleRequests) handleGuestArrives(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	bodyData := struct {
		AccompanyingGuests   int `json:"accompanying_guests"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&bodyData)
	if err != nil {
		errorResponse(w, err.Error())
		return
	}

	nameResponse(hr.guestBook.Arrived(name, bodyData.AccompanyingGuests), name, w, http.StatusAccepted, "Unable to accept guest: " + name)
}

func (hr handleRequests) handleGuestRemove(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	nameResponse(hr.guestBook.Remove(name), name, w, http.StatusAccepted, "Unable to remove guest: " + name)
}

func (hr handleRequests) handleEmptySeats(w http.ResponseWriter, r *http.Request) {
	seats := hr.guestBook.EmptySeats()
	emptySeats := struct {
		EmptySeats int `json:"emptySeats"`
	}{seats}

	outputData(emptySeats, w, http.StatusOK)
}

func nameResponse(condition bool, name string, w http.ResponseWriter, status int, errorMsg string) {
	if condition {
		response := struct {
			Name string `json:"name"`
		}{name}
		outputData(response, w, status)
	} else {
		errorResponse(w, errorMsg)
	}
}

func errorResponse(w http.ResponseWriter, errorMsg string) {
	errorData := struct {
		ErrorMsg string `json:"error"`
	}{
		errorMsg,
	}
	outputData(errorData, w, http.StatusBadRequest)
}

func outputData[T any](result T, w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Printf("Error encoding json: %v\n", err)
	}
}

func createDbConnection(settings Config.Settings) *sql.DB {
	cfg := mysql.Config{
		User:   settings.DatabaseUsername,
		Passwd: settings.DatabasePassword,
		Net:    "tcp",
		Addr:  	settings.DatabaseHost + ":" + settings.DatabasePort,
		DBName: settings.DatabaseName,
		AllowNativePasswords: true,
	}

	conn, err := sql.Open(settings.DatabaseDriver, cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := conn.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	return conn
}

func InitApp(settings Config.Settings) *mux.Router {
	dbManager := DataManager.NewDBManager(createDbConnection(settings))
	guestBook := GuestBook.NewGuestBook(dbManager, settings.VenueUsed)

	return NewHandleRequests(guestBook, settings.AppPort)
}

func main() {
	settings := Config.NewSettings()
	router := InitApp(settings)

	fmt.Println("Starting connection on port ", settings.AppPort)
	log.Fatal(http.ListenAndServe(":" + settings.AppPort, router))
}

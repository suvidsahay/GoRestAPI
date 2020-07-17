package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/joho/godotenv"
	"github.com/suvidsahay/Factly/controllers"
	"github.com/suvidsahay/Factly/types"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var a controllers.App

func TestMain(m *testing.M) {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	a.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func addUsers(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO userdb(name) VALUES($1)", "User "+strconv.Itoa(i))
	}
}

func ensureTableExists() {
	if _, err := a.DB.Exec("CREATE TABLE IF NOT EXISTS userdb (userid SERIAL, name text NOT NULL, PRIMARY KEY(userid))"); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM userdb")
	a.DB.Exec("ALTER SEQUENCE userdb_userid_seq RESTART WITH 1")
}

func getOriginalUser (id int) (types.User, error) {
	rows, err := a.DB.Query("SELECT * FROM userdb WHERE userid = $1", id)
	if err != nil {
		log.Fatal(err)
	}

	var user types.User

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Name)
		if err != nil {
			return types.User{}, err
		}
		return user, nil
	}
	return types.User{}, errors.New("no data found")
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/users", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	body := response.Body.String()
	if body != "[]\n" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestCreateUser(t *testing.T) {

	clearTable()

	var jsonStr = []byte(`{"name":"test user"}`)
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m types.User
	json.Unmarshal(response.Body.Bytes(), &m)

	if m.Name != "test user" {
		t.Errorf("Expected user name to be 'test user'. Got '%v'", m.Name)
	}

	if m.ID != 1.0 {
		t.Errorf("Expected user ID to be '1'. Got '%v'", m.ID)
	}
}

func TestUpdateUser(t *testing.T) {

	clearTable()
	addUsers(1)

	var jsonStr = []byte(`{"name":"updated name"}`)
	req, _ := http.NewRequest("PUT", "/user/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m types.User
	json.Unmarshal(response.Body.Bytes(), &m)

	var originalUser types.User

	originalUser, err := getOriginalUser(1)
	if err != nil {
		log.Fatal(err)
	}

	if m.ID != originalUser.ID {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalUser.ID, m.ID)
	}

	if m.Name != originalUser.Name {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalUser.Name, m.Name, m.Name)
	}
}

func TestDeleteUser(t *testing.T) {
	clearTable()
	addUsers(1)

	req, _ := http.NewRequest("DELETE", "/user/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}
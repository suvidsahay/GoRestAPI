package controllers

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(host, port, user, password, dbname string) {

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	var err error
	a.DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = a.DB.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	_, err = a.DB.Exec("CREATE TABLE IF NOT EXISTS userdb (userid SERIAL, name text NOT NULL, PRIMARY KEY(userid))")

	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()

	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/users", a.getUsers).Methods("GET")
	a.Router.HandleFunc("/user", a.createUser).Methods("POST")
	a.Router.HandleFunc("/user/{id}", a.updateUser).Methods("PUT")
	a.Router.HandleFunc("/user/{id}", a.deleteUser).Methods("DELETE")
}

func (a *App) Run(addr string) {
	log.Printf("\nServer starting on port %s", addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

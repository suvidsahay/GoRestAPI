package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/suvidsahay/Factly/database"
	"github.com/suvidsahay/Factly/types"
	"github.com/suvidsahay/Factly/responses"
	"net/http"
)

func (a *App) getUsers(w http.ResponseWriter, r *http.Request) {

	users, err := database.GetUsers(a.DB)
	if err != nil {
		responses.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, users)
}

func (a *App) createUser(w http.ResponseWriter, r *http.Request) {
	var user types.User
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&user)

	if err != nil {
		responses.JSONError(w, http.StatusBadRequest, err)
		return
	}

	err = database.CreateUser(a.DB, &user)
	if err != nil {
		responses.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, user)
}

func (a *App) updateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user types.User
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&user)
	if err != nil {
		responses.JSONError(w, http.StatusBadRequest, err)
		return
	}

	user, err = database.UpdateUser(a.DB, params["id"], user.Name)
	if err != nil {
		responses.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, user)
}

func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	rows, err := database.DeleteUser(a.DB, params["id"])
	if err != nil {
		responses.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	if rows == 0 {
		responses.JSONError(w, http.StatusBadRequest, errors.New("data not found"))
		return
	}

	responses.JSON(w, http.StatusOK, `{result:success}`)
}
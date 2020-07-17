package database

import (
	"database/sql"
	"github.com/suvidsahay/Factly/types"
	"strconv"
)

var users []types.User

func GetUsers(db *sql.DB) ([]types.User, error) {
	var user types.User
	users = []types.User{}

	rows, err := db.Query("SELECT * fROM userdb")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Name)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func CreateUser(db *sql.DB, user *types.User) error {
	err := db.QueryRow("INSERT INTO userdb (name) VALUES($1) RETURNING userid;", user.Name).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}

func UpdateUser(db *sql.DB, id string, name string) (types.User, error) {
	var user types.User

	_, err := db.Exec("UPDATE userdb set name = $1 WHERE userid = $2;", name, id)
	if err != nil {
		return types.User{}, err
	}

	user.ID, err = strconv.Atoi(id)
	user.Name = name
	if err != nil {
		return types.User{}, err
	}

	return user, nil
}

func DeleteUser(db *sql.DB, id string) (int64, error) {
	result, err := db.Exec("DELETE FROM userdb where userid = $1", id)
	if err != nil {
		return 0, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, nil
}

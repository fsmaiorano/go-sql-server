package database

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Create a new user in database
func CreateUser(w http.ResponseWriter, r *http.Request) {
	request, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.Write([]byte("Request body is empty: " + err.Error()))
	}

	var user User

	if err = json.Unmarshal(request, &user); err != nil {
		w.Write([]byte("Error unmarshalling request body: " + err.Error()))
	}

	db, err := Connect()

	if err != nil {
		w.Write([]byte("Error connecting to database: " + err.Error()))
		return
	}

	defer db.Close()

	statement, err := db.Prepare("INSERT INTO usuarios (name, email) VALUES (@name, @email);")
	if err != nil {
		w.Write([]byte("Error preparing statement: " + err.Error()))
		return
	}

	defer statement.Close()

	_, err = statement.Exec(sql.Named("name", user.Name), sql.Named("email", user.Email))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error inserting user " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User inserted with success"))
}

// Find all users in database
func FindUsers(w http.ResponseWriter, r *http.Request) {
	db, err := Connect()
	if err != nil {
		w.Write([]byte("Error connecting to database: " + err.Error()))
		return
	}

	defer db.Close()

	lines, err := db.Query("SELECT * FROM usuarios;")
	if err != nil {
		w.Write([]byte("Error preparing statement: " + err.Error()))
		return
	}

	defer lines.Close()

	var users []User
	for lines.Next() {
		var user User

		if err := lines.Scan(&user.Id, &user.Name, &user.Email); err != nil {
			w.Write([]byte("Error scanning line: " + err.Error()))
			return
		}

		users = append(users, user)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		w.Write([]byte("Error encoding users: " + err.Error()))
		return
	}
}

// Find a user in database
func FindUser(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)

	ID, err := strconv.ParseUint(parameters["id"], 10, 32)
	if err != nil {
		w.Write([]byte("Error parsing id: " + err.Error()))
		return
	}

	db, err := Connect()
	if err != nil {
		w.Write([]byte("Error connecting to database: " + err.Error()))
		return
	}

	line, err := db.Query("SELECT * FROM usuarios WHERE id = @id;", sql.Named("id", ID))
	if err != nil {
		w.Write([]byte("Error preparing statement: " + err.Error()))
		return
	}

	var user User
	if line.Next() {
		if err := line.Scan(&user.Id, &user.Name, &user.Email); err != nil {
			w.Write([]byte("Error scanning line: " + err.Error()))
			return
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		w.Write([]byte("Error encoding user: " + err.Error()))
		return
	}
}

// Update a user in database
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)

	ID, err := strconv.ParseUint(parameters["id"], 10, 32)
	if err != nil {
		w.Write([]byte("Error parsing id: " + err.Error()))
		return
	}

	request, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.Write([]byte("Request body is empty: " + err.Error()))
	}

	var user User

	if err = json.Unmarshal(request, &user); err != nil {
		w.Write([]byte("Error unmarshalling request body: " + err.Error()))
	}

	db, err := Connect()

	if err != nil {
		w.Write([]byte("Error connecting to database: " + err.Error()))
		return
	}

	defer db.Close()

	statement, err := db.Prepare("UPDATE usuarios SET name = @name, email = @email WHERE id = @id;")
	if err != nil {
		w.Write([]byte("Error preparing statement: " + err.Error()))
		return
	}

	defer statement.Close()

	_, err = statement.Exec(sql.Named("id", ID), sql.Named("name", user.Name), sql.Named("email", user.Email))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error updating user " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User updated with success"))
}

// Delete a user in database
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)

	ID, err := strconv.ParseUint(parameters["id"], 10, 32)
	if err != nil {
		w.Write([]byte("Error parsing id: " + err.Error()))
		return
	}

	db, err := Connect()

	if err != nil {
		w.Write([]byte("Error connecting to database: " + err.Error()))
		return
	}

	defer db.Close()

	statement, err := db.Prepare("DELETE FROM usuarios WHERE id = @id;")
	if err != nil {
		w.Write([]byte("Error preparing statement: " + err.Error()))
		return
	}

	defer statement.Close()

	_, err = statement.Exec(sql.Named("id", ID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error deleting user " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User deleted with success"))
}

package router

import (
	"net/http"
	"sql-server/src/database"

	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/users", database.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/users", database.FindUsers).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", database.FindUser).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", database.UpdateUser).Methods(http.MethodPut)
	router.HandleFunc("/users/{id}", database.DeleteUser).Methods(http.MethodDelete)

	return router
}

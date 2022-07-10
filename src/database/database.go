package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/denisenkom/go-mssqldb/azuread"
)

func Connect() (*sql.DB, error) {

	fmt.Println(azuread.DriverName)

	db, err := sql.Open("sqlserver", "sqlserver://sa:Password123@localhost:1433?database=devbook")

	if err != nil {
		fmt.Println("ERROR!!!")
		log.Fatal(err.Error())
		return nil, err
	}

	if err = db.Ping(); err != nil {
		fmt.Println("ERROR!!! PING")
		log.Fatal(err.Error())
		return nil, err
	}

	return db, nil
}

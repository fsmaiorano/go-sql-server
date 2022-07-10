package main

import (
	"fmt"
	"log"
	"net/http"
	"sql-server/src/router"
)

func main() {

	router := router.InitRoutes()

	fmt.Println("Listening on port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}

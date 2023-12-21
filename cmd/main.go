package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kamalkish0r/GO-REST/pkg/db"
	"github.com/kamalkish0r/GO-REST/pkg/routes"
)

func main() {
	// connect to db
	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// create tables if they don't exist
	err = db.CreateTables(database)
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	routes.Routes(router)

	log.Println("Started server at port : ", 8080)
	http.ListenAndServe(":8080", router)
}

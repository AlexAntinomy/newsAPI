package main

import (
	"log"
	"net/http"
	"news/pkg/api"
	"news/pkg/storage"
	"news/pkg/storage/postgres"
)

func main() {
	var srv struct {
		db  storage.Interface
		api *api.API
	}

	var db storage.Interface
	var err error

	// PostgreSQL (удалённый сервер).
	db, err = postgres.New("postgres://postgres:123@87.242.119.32:5432/news?sslmode=disable")
	if err != nil {
		log.Fatal("Postgres connection error:", err)
	}

	// db, err = mongo.New("mongodb://87.242.119.32:27017")
	// if err != nil {
	// 	log.Fatal("Mongo connection error:", err)
	// }

	// БД в памяти
	// db := memdb.New()

	srv.db = db
	srv.api = api.New(srv.db)

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", srv.api.Router()); err != nil {
		log.Fatal("Server error:", err)
	}
}

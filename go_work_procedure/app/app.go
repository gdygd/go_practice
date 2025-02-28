package app

import (
	"db"
	"net/http"

	"github.com/gorilla/mux"
)

type AppHandler struct {
	http.Handler
	dbHandler *db.DBHandler
}

func MakeAppHandler() AppHandler {
	r := mux.NewRouter()

	// db information
	newDBHandler := &db.DBHandler{
		Host: "10.1.0.115", Port: 3306, DBname: "test", User: "dev", Password: "dev",
	}

	a := AppHandler{
		Handler:   r,
		dbHandler: newDBHandler,
	}

	r.HandleFunc("/", a.IndexHandler).Methods("GET")
	r.HandleFunc("/write", a.WriteHandler).Methods("POST")
	r.HandleFunc("/read", a.ReadHandler).Methods("GET")

	return a
}

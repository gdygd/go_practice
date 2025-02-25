package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func InitializeRoute() *mux.Router {
	log.Print("InitializeRoute")

	router := mux.NewRouter()
	router.HandleFunc("/signup", SignUp).Methods("POST")
	router.HandleFunc("/signin", SignIn).Methods("POST")
	router.HandleFunc("/apitest", Apitest).Methods("POST")

	// middleware
	router.Use(AuthMiddleware)

	return router
}

func main() {
	log.Print("program start")
	router := InitializeRoute()

	log.Print("program run..")

	log.Fatal(http.ListenAndServe(":3002", handlers.LoggingHandler(os.Stdout, router)))
}

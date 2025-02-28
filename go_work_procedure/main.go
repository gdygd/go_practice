package main

import (
	"app"
	"net/http"
)

func main() {
	a := app.MakeAppHandler()
	http.ListenAndServe(":5000", a)
}

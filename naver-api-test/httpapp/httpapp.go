package httpapp

import (
	"net/http"

	"github.com/gorilla/mux"
)

//------------------------------------------------------------------------------
// HttpAppHandler
//------------------------------------------------------------------------------
type HttpAppHandler struct {
	http.Handler
}

//------------------------------------------------------------------------------
// MakeHandler
//------------------------------------------------------------------------------
func MakeHandler() *HttpAppHandler {

	r := mux.NewRouter().StrictSlash(true)
	a := &HttpAppHandler{
		Handler: r,
	}

	// Init API
	r.HandleFunc("/naverapi", a.GetNaverApi).Methods("GET")

	return a
}

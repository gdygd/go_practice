package app

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (a AppHandler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "[Index] Hello World")
}

func (a *AppHandler) ReadHandler(w http.ResponseWriter, r *http.Request) {
	res := &Result{Res: 1}

	sliceInput, err := a.dbHandler.ReadInput()
	if err != nil {
		res.Msg = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	res.Res = 0
	res.Msg = "success"
	res.Data = sliceInput
	json.NewEncoder(w).Encode(res)
}

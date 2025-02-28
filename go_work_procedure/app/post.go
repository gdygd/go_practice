package app

import (
	"db"
	"encoding/json"
	"io"
	"net/http"
)

func (a AppHandler) WriteHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	res := &Result{Res: 1}

	if err != nil {
		res.Msg = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	insInp := db.InsertInput{}
	err = json.Unmarshal(body, &insInp)
	if err != nil {
		res.Msg = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	// DB 가서 insert
	err = a.dbHandler.WriteInput(insInp.Str)
	if err != nil {
		res.Msg = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	res.Res = 0
	res.Msg = "success"
	json.NewEncoder(w).Encode(res)
}

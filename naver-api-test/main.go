package main

import (
	"fmt"
	"naver-api-test/httpapp"
	"net/http"
)

var httpInst *httpapp.HttpAppHandler = nil

func main() {
	fmt.Printf("naver api test start..")

	httpInst = httpapp.MakeHandler()

	//err := http.ListenAndServe(":5500", handlers.CORS(originsOK, headersOK, methodsOK)(appHnd.Handler))
	err := http.ListenAndServe(":7500", httpInst.Handler)

	if err != nil {
		panic(err)
	}

	// http.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
	// 	w.Write([]byte("Hello World"))
	// })

	// http.ListenAndServe(":7500", nil)
}

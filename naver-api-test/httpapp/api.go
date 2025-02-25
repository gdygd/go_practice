package httpapp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//------------------------------------------------------------------------------
// GetNaverApi post
//------------------------------------------------------------------------------
func (a *HttpAppHandler) GetNaverApi(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("GetNaverApi, /naverapi")

	apiUrl := "https://naveropenapi.apigw.ntruss.com/map-direction-15/v1/driving?start=126.974710,37.344092&goal=126.96976883781271,37.358512374655966&option=trafast"

	fmt.Printf("GetNaverApi, /naverapi  requrl : %v", apiUrl)

	// Request 객체 생성
	//req, err := http.NewRequest("GET", "http://csharp.tips/feed/rss", nil)
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		panic(err)
	}

	//필요시 헤더 추가 가능
	req.Header.Add("User-Agent", "Crawler")

	// Client객체에서 Request 실행
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 결과 출력
	bytes, _ := ioutil.ReadAll(resp.Body)
	str := string(bytes) //바이트를 문자열로
	fmt.Println(str)
	//returnData.Result = result
	json.NewEncoder(w).Encode(str)
}

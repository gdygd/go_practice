package main

import (
	"fmt"
	"golang.org/x/sync/singleflight"
)

var g singleflight.Group
var cache = map[string]string{} // 간단한 캐시 흉내

func getUserFromDB(userID string) (interface{}, error) {
	// 실제로는 DB 쿼리
	fmt.Println("DB 조회:", userID)
	return "user-data-for-" + userID, nil
}

func GetUser(userID string) (string, error) {
	if v, ok := cache[userID]; ok {
		return v, nil // 캐시에 있으면 바로 반환
	}

	// 캐시에 없으면 singleflight로 DB 조회 (동시 요청 시 1번만 조회)
	v, err, _ := g.Do(userID, func() (interface{}, error) {
		return getUserFromDB(userID)
	})
	if err != nil {
		return "", err
	}

	result := v.(string)
	cache[userID] = result
	return result, nil
}

func main() {
	// 여러 요청이 동시에 같은 유저를 조회해도 DB는 1번만 호출됨
	for i := 0; i < 3; i++ {
		go GetUser("user-1")
	}
}
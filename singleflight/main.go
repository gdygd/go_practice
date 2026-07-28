package main

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

func main() {
	var g singleflight.Group
	var wg sync.WaitGroup

	// 실제로는 "무거운 작업" (예: DB 조회, 외부 API 호출)
	expensiveWork := func() (interface{}, error) {
		fmt.Println("실제 작업 실행!") // 이게 몇 번 찍히는지 확인해보세요
		time.Sleep(1 * time.Second)
		return "결과값", nil
	}

	// 동시에 5개의 goroutine이 같은 key로 요청
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			v, err, shared := g.Do("same-key", expensiveWork)
			fmt.Printf("goroutine %d: 결과=%v, 공유됨=%v, err=%v\n", id, v, shared, err)
		}(i)
	}

	time.Sleep(time.Second * 1)
	v, err, shared := g.Do("same-key", expensiveWork)
	fmt.Printf("##goroutine %d: 결과=%v, 공유됨=%v, err=%v\n", 6, v, shared, err)

	wg.Wait()
}

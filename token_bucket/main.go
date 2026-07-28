package main

import (
	"net/http"

	"token_bucket/tokenbucket"

	"github.com/gin-gonic/gin"
)

// ResourceManage
type ResourceManage struct {
	ciriticalBucket   *tokenbucket.TokenBucket
	normalBucket      *tokenbucket.TokenBucket
	lowPriorityBucket *tokenbucket.TokenBucket
}

func NewResourceManager() *ResourceManage {
	return &ResourceManage{
		// 많은 토큰과 버스트 용량
		ciriticalBucket: tokenbucket.NewTokenBucket(100, 50),
		// 일반 작업에 적당한 토큰과 버스트 용량
		normalBucket: tokenbucket.NewTokenBucket(50, 20),
		// 낮은우선순위 작업에 적당한 토큰과 버스트 용량
		lowPriorityBucket: tokenbucket.NewTokenBucket(20, 5),
	}
}

var (
	tb     *tokenbucket.TokenBucket = nil
	resMng *ResourceManage          = nil
)

func testapi(ctx *gin.Context) {
	if !tb.Take(1) {
		ctx.JSON(http.StatusTooManyRequests, "too many request...")
		return
	}

	// if !tb.TakeWithTimeout(1, time.Second*1) {
	// 	ctx.JSON(http.StatusTooManyRequests, "too many request...")
	// 	return
	// }

	ctx.JSON(http.StatusOK, "hello")
}

func testapi2(ctx *gin.Context) {
	if !resMng.lowPriorityBucket.Take(1) {
		ctx.JSON(http.StatusTooManyRequests, "too many request...")
		return
	}

	// if !tb.TakeWithTimeout(1, time.Second*1) {
	// 	ctx.JSON(http.StatusTooManyRequests, "too many request...")
	// 	return
	// }

	ctx.JSON(http.StatusOK, "hello")
}

func main() {
	tb = tokenbucket.NewTokenBucket(5, 5) // bucket size 10, rate 5 per a second
	resMng = NewResourceManager()

	r := gin.Default()

	r.GET("/tbtest", testapi)   // tocken bucket test api
	r.GET("/tbtest2", testapi2) // tocken bucket test api by resource manager

	// 서버 실행
	r.Run(":9100")
}

package main

import (
	"net/http"

	"token_bucket/tokenbucket"

	"github.com/gin-gonic/gin"
)

var tb *tokenbucket.TokenBucket = nil

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

func main() {
	tb = tokenbucket.NewTokenBucket(5, 5) // bucket size 10, rate 5 per a second
	r := gin.Default()

	r.GET("/tbtest", testapi) // tocken bucket test api

	// 서버 실행
	r.Run(":9100")
}

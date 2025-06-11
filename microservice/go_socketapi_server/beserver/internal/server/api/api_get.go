package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (server *Server) testapi(ctx *gin.Context) {

	time.Sleep(time.Microsecond * 3000)
	ctx.JSON(http.StatusOK, "hello")
}

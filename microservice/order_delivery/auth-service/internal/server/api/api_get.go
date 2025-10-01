package api

import (
	"auth-service/internal/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (server *Server) testapi(ctx *gin.Context) {

	time.Sleep(time.Microsecond * 3000)

	strdt, err := server.dbHnd.ReadSysdate(ctx)
	if err != nil {
		logger.Log.Error("testapi err..%v", err)
	}
	logger.Log.Print(2, "testapi :%v", strdt)

	ctx.JSON(http.StatusOK, "hello")

}

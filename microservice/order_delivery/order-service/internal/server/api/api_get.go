package api

import (
	"net/http"
	"order-service/internal/logger"
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

func (server *Server) getOrderInfo(ctx *gin.Context) {
	var req orderInfoRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	orders, err := server.dbHnd.ReadOrderInfo(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := []OrderResponse{}
	for _, ord := range orders {
		rsp = append(rsp, convertOrder(ord))
	}

	ctx.JSON(http.StatusOK, rsp)
}

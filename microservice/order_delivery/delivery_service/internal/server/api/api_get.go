package api

import (
	"delivery_service/internal/logger"
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

func (server *Server) getDeliveryInfo(ctx *gin.Context) {
	var req deliveryInfoRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	deliveries, err := server.dbHnd.ReadDeliveries(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := []DeliveryResponse{}
	for _, deli := range deliveries {
		rsp = append(rsp, convertDelivery(deli))
	}

	ctx.JSON(http.StatusOK, rsp)
}

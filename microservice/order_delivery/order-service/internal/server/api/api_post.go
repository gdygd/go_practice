package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/internal/logger"

	"github.com/gdygd/goglib"
	"github.com/gin-gonic/gin"
)

func (server *Server) requestOrder(ctx *gin.Context) {
	var req orderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ord, err := server.dbHnd.RequestOrder(ctx, getOrderPrarm(req))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, convertOrder(ord))
}

const BASE_URL = "http://10.1.0.119:9080"

func (server *Server) requestOrder2(ctx *gin.Context) {
	var req orderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ord, err := server.dbHnd.RequestOrder(ctx, getOrderPrarm(req))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// request orchestrator
	// payload, _ := json.Marshal(ord)
	// url := "http://10.1.0.119:9080/saga/order"
	// resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("orchestrator 호출 실패: %w", err)))
	// }
	// defer resp.Body.Close()

	logger.Log.Print(2, "Request saga order ... ")

	sagaObj := convertSagaOrder(ord)
	payload, _ := json.Marshal(sagaObj)
	logger.Log.Print(2, "id:%d, nm:%s, amount:%d", sagaObj.OrderId, sagaObj.Username, sagaObj.TotalAmout)

	url := fmt.Sprintf("%s/saga/order", BASE_URL)
	statuscode, _, err := goglib.HttpRequest(ctx, payload, "POST", url)
	if err != nil {
		logger.Log.Print(2, "Request saga order fail... ")
		ctx.JSON(statuscode, errorResponse(fmt.Errorf("orchestrator 호출 실패: %w", err)))
		return
	}

	ctx.JSON(http.StatusOK, convertOrder(ord))
}

type orderStateRequest struct {
	ID int `uri:"id" binding:"required,min=1"`
}

func (server *Server) cancelOrder(ctx *gin.Context) {
	var req orderStateRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.dbHnd.CancelOrder(ctx, int(req.ID))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func (server *Server) confirmOrder(ctx *gin.Context) {
	var req orderStateRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.dbHnd.ConfirmOrder(ctx, int(req.ID))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

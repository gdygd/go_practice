package api

import (
	"fmt"
	"net/http"
	"order-delivery-saga/internal/logger"

	"github.com/gin-gonic/gin"
)

func (server *Server) sagaOrder(ctx *gin.Context) {
	var req orderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.Print(2, "body parsing error.. %v", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	logger.Log.Print(2, "orchestrator start, order : %d", req.OrderId)

	// step 1 : call delivery service
	statecode, err := callDeliveryService(ctx, req.OrderId)
	if err != nil {
		logger.Log.Print(2, "call delivery service error.. %v", err)
		ctx.JSON(statecode, errorResponse(fmt.Errorf("delivery 호출 실패: %w", err)))
		return
	}

	// confirm order

	logger.Log.Print(2, "orchestrator end, order : %d", req.OrderId)

	ctx.JSON(http.StatusOK, nil)
}

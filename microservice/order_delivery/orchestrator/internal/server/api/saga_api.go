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
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	logger.Log.Print(2, "orchestrator start, order : %d", req.ORDER_ID)

	// step 1 : call delivery service
	if err := callDeliveryService(req.ORDER_ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("delivery 호출 실패: %w", err)))
	}

	// confirm order

	logger.Log.Print(2, "orchestrator end, order : %d", req.ORDER_ID)

	ctx.JSON(http.StatusOK, nil)
}

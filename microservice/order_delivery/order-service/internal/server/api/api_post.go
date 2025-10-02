package api

import (
	"net/http"

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

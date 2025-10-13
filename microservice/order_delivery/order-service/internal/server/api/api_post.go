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

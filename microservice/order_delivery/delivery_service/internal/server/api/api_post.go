package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) requestDelivery(ctx *gin.Context) {
	var req deliveryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	deli, err := server.dbHnd.RequestDelivery(ctx, getDeliveryPrarm(req))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, convertDelivery(deli))
}

type deliveryStateRequest struct {
	ID int `uri:"id" binding:"required,min=1"`
}

func (server *Server) cancelDelivery(ctx *gin.Context) {
	var req deliveryStateRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.dbHnd.CancelDelivery(ctx, int(req.ID))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func (server *Server) confirmDelivery(ctx *gin.Context) {
	var req deliveryStateRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.dbHnd.ConfirmDelivery(ctx, int(req.ID))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

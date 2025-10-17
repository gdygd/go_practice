package api

import (
	"auth-service/internal/db"
	"auth-service/internal/logger"
	"auth-service/internal/util"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func (server *Server) userLogin(ctx *gin.Context) {

	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	logger.Log.Print(2, "userLogin..#1 %v", req)

	user, err := server.dbHnd.ReadUser(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	logger.Log.Print(2, "userLogin..#2 %v", user)

	err = util.CheckPassword(req.Password, user.PASSWD)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	logger.Log.Print(2, "userLogin..check.. %v", user)

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.USER_NM,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	logger.Log.Print(2, "userLogin..acc token %v", accessToken)

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.USER_NM,
		server.config.RefreshTokenDuration,
	)

	logger.Log.Print(2, "userLogin..ref token %v", refreshToken)

	ssparam := db.SESSIONS{
		ID:         refreshPayload.ID.String(),
		USER_NM:    user.USER_NM,
		REF_TOKEN:  refreshToken,
		USER_AGENT: ctx.Request.UserAgent(),
		CLIENT_IP:  ctx.ClientIP(),
		BLOCK_YN:   0,
		EXP_DT:     refreshPayload.ExpiredAt,
	}

	se, err := server.dbHnd.CreateSession(ctx, ssparam)

	seid, _ := uuid.Parse(se.ID)
	rsp := loginUserResponse{
		SessionID:             seid,
		AcessToken:            accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) tokenVerify(ctx *gin.Context) {

	authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
	if len(authorizationHeader) == 0 {
		err := errors.New("authorization header is not provided")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	fields := strings.Fields(authorizationHeader)
	if len(fields) < 2 {
		err := errors.New("invalid authorization header format")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != authorizationTypeBearer {
		err := fmt.Errorf("unsupported authorization type %s", authorizationType)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken := fields[1]
	payload, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, payload)
}

func (server *Server) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	se, err := server.dbHnd.ReadSession(ctx, refreshPayload.ID.String())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// check session
	if se.BLOCK_YN == 1 {
		err := fmt.Errorf("session is blocked..")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// chek user
	if se.USER_NM != refreshPayload.Username {
		err := fmt.Errorf("incorrect session user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// check token
	if se.REF_TOKEN != req.RefreshToken {
		err := fmt.Errorf("mismatched session token")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// check exp date
	if time.Now().After(se.EXP_DT) {
		err := fmt.Errorf("expired session token")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		refreshPayload.Username,
		server.config.AccessTokenDuration,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := renewAccessTokenResponse{
		AcessToken:           accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}

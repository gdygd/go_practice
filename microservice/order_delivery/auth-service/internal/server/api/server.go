package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"auth-service/internal/config"
	"auth-service/internal/container"
	"auth-service/internal/db"
	"auth-service/internal/logger"
	"auth-service/internal/memory"
	"auth-service/internal/service"
	apiserv "auth-service/internal/service/api"

	"github.com/gdygd/goglib/token"

	"github.com/gin-gonic/gin"
)

var TokenSecretKey = "asdFQWER!@#$ASDFEWR#@$~!~!@#123"

// Server serves HTTP requests for our banking service.
type Server struct {
	wg         *sync.WaitGroup
	srv        *http.Server
	config     *config.Config
	tokenMaker token.Maker
	router     *gin.Engine
	service    service.ServiceInterface
	dbHnd      db.DbHandler
	objdb      *memory.RedisDb
}

func NewServer(wg *sync.WaitGroup, ct *container.Container) (*Server, error) {
	// init service
	apiservice := apiserv.NewApiService(ct.DbHnd, ct.ObjDb)
	tokenMaker, err := token.NewJWTMaker(TokenSecretKey)

	if err != nil {
		return nil, fmt.Errorf("cannot create token maker:%w", err)
	}

	server := &Server{
		wg:         wg,
		config:     ct.Config,
		tokenMaker: tokenMaker,
		service:    apiservice,
		dbHnd:      ct.DbHnd,
		objdb:      ct.ObjDb,
	}

	server.setupRouter()

	server.srv = &http.Server{}
	server.srv.Addr = ct.Config.HTTPServerAddress
	server.srv.Handler = server.router.Handler()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	// gin.SetMode(gin.DebugMode)
	fmt.Printf("%v, \n", server.config.AllowOrigins)

	addresses := strings.Split(server.config.AllowOrigins, ",")

	router.Use(corsMiddleware(addresses))
	router.Use(authMiddleware(server.tokenMaker))
	router.GET("/test", server.testapi)
	router.POST("/auth/login", server.userLogin)
	router.POST("/auth/verify", server.tokenVerify)
	router.POST("/auth/refresh", server.renewAccessToken)

	server.router = router
}

func (server *Server) Start() error {
	logger.Log.Print(2, "Gin server start.")

	if err := server.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Log.Error("listen error. %v", err)
		return err
	}

	return nil
}

func (server *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer server.wg.Done()
	if err := server.srv.Shutdown(ctx); err != nil {
		logger.Log.Error("Server Shutdown:", err)
		return err
	}
	return nil
}

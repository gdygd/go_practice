package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"order-service/internal/config"
	"order-service/internal/container"
	"order-service/internal/db"
	"order-service/internal/logger"
	"order-service/internal/memory"
	"order-service/internal/service"
	apiserv "order-service/internal/service/api"

	"github.com/gdygd/goglib/token"

	"github.com/gin-gonic/gin"
)

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
	tokenMaker, err := token.NewJWTMaker(ct.Config.TokenSecretKey)

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
	router.GET("/info", server.getOrderInfo)
	router.POST("/request", server.requestOrder)
	router.POST("/cancel", server.cancelOrder)

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

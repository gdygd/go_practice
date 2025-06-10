package api

import (
	"context"
	"fmt"
	"net/http"
	"server/config"
	"server/restserver/internal/container"
	"server/restserver/internal/logger"
	"server/restserver/internal/service"
	apiserv "server/restserver/internal/service/api"
	"server/token"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	srv        *http.Server
	config     *config.Config
	tokenMaker token.Maker
	router     *gin.Engine
	service    service.ApiServiceInterface
}

func NewServer(ct *container.Container) (*Server, error) {

	// init service
	apiservice := apiserv.NewApiService(ct.DbHnd, ct.ObjDb)

	tokenMaker, err := token.NewJWTMaker("aaaaa")
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     ct.Config,
		tokenMaker: tokenMaker,
		service:    apiservice,
	}

	server.setupRouter()

	server.srv = &http.Server{}
	server.srv.Addr = "0.0.0.0:9090"
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

	server.router = router

}

func (server *Server) Start() error {
	logger.Apilog.Print(2, "Gin server start.")

	if err := server.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Apilog.Error("listen error. %v", err)
		return err
	}

	return nil
}

func (server *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.srv.Shutdown(ctx); err != nil {
		logger.Apilog.Error("Server Shutdown:", err)
		return err
	}
	return nil
}

// func resolveAddress(addr ...string) string {
// 	switch len(addr) {
// 	case 0:
// 		if port := os.Getenv("PORT"); port != "" {
// 			log.Printf("Environment variable PORT=\"%s\"", port)
// 			return ":" + port
// 		}
// 		log.Println("Environment variable PORT is undefined. Using port :8080 by default")
// 		return ":8080"
// 	case 1:
// 		return addr[0]
// 	default:
// 		panic("too many parameters")
// 	}
// }

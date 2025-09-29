package api

import (
	"net/http"
	"sync"

	"auth-service/config"
	"auth-service/internal/container"
	"auth-service/token"

	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	wg         *sync.WaitGroup
	srv        *http.Server
	config     *config.Config
	tokenMaker token.Maker
	router     *gin.Engine
	service    service.ApiServiceInterface
}

func NewServer(wg *sync.WaitGroup, ct *container.Container) (*Server, error) {

	return nil, nil
}

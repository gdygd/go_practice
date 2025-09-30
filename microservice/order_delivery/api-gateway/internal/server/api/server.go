package api

import (
	"api-gateway/internal/container"
	"net/http"
	"sync"

	"github.com/gdygd/goglib/token"

	"github.com/gin-gonic/gin"
	"github.com/jcmturner/gokrb5/v8/config"
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

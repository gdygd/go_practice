package gapi

import (
	"auth-service/internal/config"
	"auth-service/internal/container"
	"auth-service/internal/db"
	"auth-service/internal/logger"
	"auth-service/internal/memory"
	"auth-service/internal/service"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	apiserv "auth-service/internal/service/api"
	pb "auth-service/pb/proto"

	"github.com/gdygd/goglib/token"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	wg *sync.WaitGroup
	pb.UnimplementedAuthServiceServer
	gServer *grpc.Server

	config       *config.Config
	tokenMaker   token.Maker
	service      service.ServiceInterface
	dbHnd        db.DbHandler
	objdb        *memory.RedisDb
	ch_terminate chan bool
}

func NewServer(wg *sync.WaitGroup, ct *container.Container, ch_terminate chan bool) (*Server, error) {
	apiservice := apiserv.NewApiService(ct.DbHnd, ct.ObjDb)
	tokenMaker, err := token.NewJWTMaker(ct.Config.TokenSecretKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker:%w", err)
	}

	server := &Server{
		wg:           wg,
		config:       ct.Config,
		tokenMaker:   tokenMaker,
		service:      apiservice,
		dbHnd:        ct.DbHnd,
		objdb:        ct.ObjDb,
		ch_terminate: ch_terminate,
	}

	grpcLogger := grpc.UnaryInterceptor(GrpcServerLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterAuthServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	return server, nil
}

func NewGatewayServer(wg *sync.WaitGroup, ct *container.Container, ch_terminate chan bool) (*Server, error) {
	return nil, nil
}

func (server *Server) StartgPRC() error {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Log.Error("cannot create listener:")
	}

	err = server.gServer.Serve(listener)
	if err != nil {
		if errors.Is(err, grpc.ErrServerStopped) {
			return nil
		}
		logger.Log.Error("gRPC server faield to serve, err:%v", err)
		return err
	}
	return nil
}

func (server *Server) Shutdown() error {
	defer server.wg.Done()
	server.gServer.GracefulStop()

	return nil
}

func (server *Server) setupRouter() {
}

func (server *Server) testapi(ctx *gin.Context) {
	time.Sleep(time.Microsecond * 3000)

	strdt, err := server.dbHnd.ReadSysdate(ctx)
	if err != nil {
		logger.Log.Error("testapi err..%v", err)
	}
	logger.Log.Print(2, "testapi :%v", strdt)

	ctx.JSON(http.StatusOK, "hello")
}

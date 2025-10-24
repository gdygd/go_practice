package gapi

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"grpc_svr_test/internal/config"
	"grpc_svr_test/internal/container"
	"grpc_svr_test/internal/db"
	"grpc_svr_test/internal/logger"
	"grpc_svr_test/internal/memory"
	"grpc_svr_test/internal/service"

	apiserv "grpc_svr_test/internal/service/api"
	pb "grpc_svr_test/pb"

	"github.com/gdygd/goglib/token"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	R_TIME_OUT = 5 * time.Second
	W_TIME_OUT = 5 * time.Second
)

type Server struct {
	wg *sync.WaitGroup
	pb.UnimplementedHelloServiceServer
	gServer      *grpc.Server
	rServer      *http.Server
	router       *gin.Engine
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

	// grpcLogger := grpc.UnaryInterceptor(GrpcServerLogger)
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			GrpcServerLogger,
		),
	)
	pb.RegisterHelloServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	server.gServer = grpcServer

	return server, nil
}

// func NewGatewayServer(wg *sync.WaitGroup, ct *container.Container, ch_terminate chan bool) (*Server, error) {
// 	apiservice := apiserv.NewApiService(ct.DbHnd, ct.ObjDb)
// 	tokenMaker, err := token.NewJWTMaker(ct.Config.TokenSecretKey)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot create token maker:%w", err)
// 	}

// 	server := &Server{
// 		wg:           wg,
// 		config:       ct.Config,
// 		tokenMaker:   tokenMaker,
// 		service:      apiservice,
// 		dbHnd:        ct.DbHnd,
// 		objdb:        ct.ObjDb,
// 		ch_terminate: ch_terminate,
// 	}

// 	server.rServer = &http.Server{}
// 	server.rServer.Addr = "0.0.0.0:9091"

// 	server.rServer.ReadTimeout = R_TIME_OUT
// 	server.rServer.WriteTimeout = W_TIME_OUT

// 	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
// 		MarshalOptions: protojson.MarshalOptions{
// 			UseProtoNames: true,
// 		},
// 		UnmarshalOptions: protojson.UnmarshalOptions{
// 			DiscardUnknown: true,
// 		},
// 	})

// 	ctx := context.Background()
// 	// ctx, cancel := context.WithCancel(ctx)

// 	grpcMux := runtime.NewServeMux(jsonOption)
// 	opts := []grpc.DialOption{grpc.WithInsecure()}
// 	err = pb.RegisterAuthServiceHandlerFromEndpoint(
// 		ctx, grpcMux, "localhost:9090", opts,
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot register grpc gateway handler: %w", err)
// 	}

// 	server.setupRouter()
// 	server.router.Any("/v1/*any", gin.WrapH(grpcMux))
// 	server.rServer.Handler = server.router.Handler()

// 	return server, nil
// }

// func NewGatewayServer(wg *sync.WaitGroup, ct *container.Container, ch_terminate chan bool) (*Server, error) {
// 	apiservice := apiserv.NewApiService(ct.DbHnd, ct.ObjDb)
// 	tokenMaker, err := token.NewJWTMaker(ct.Config.TokenSecretKey)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot create token maker:%w", err)
// 	}

// 	server := &Server{
// 		wg:           wg,
// 		config:       ct.Config,
// 		tokenMaker:   tokenMaker,
// 		service:      apiservice,
// 		dbHnd:        ct.DbHnd,
// 		objdb:        ct.ObjDb,
// 		ch_terminate: ch_terminate,
// 	}

// 	router := gin.Default()

// 	// gRPC-Gateway Mux 생성
// 	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
// 		MarshalOptions: protojson.MarshalOptions{
// 			UseProtoNames: true,
// 		},
// 		UnmarshalOptions: protojson.UnmarshalOptions{
// 			DiscardUnknown: true,
// 		},
// 	})

// 	ctx := context.Background()
// 	// ctx, cancel := context.WithCancel(ctx)

// 	grpcMux := runtime.NewServeMux(jsonOption)

// 	// 실제 gRPC 서버 localhost:50051)에 연결
// 	opts := []grpc.DialOption{grpc.WithInsecure()}
// 	err = pb.RegisterAuthServiceHandlerFromEndpoint(
// 		// ctx, grpcMux, "localhost:50051", opts,
// 		ctx, grpcMux, server.config.GRPCServerAddress, opts,
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot register grpc gateway handler: %w", err)
// 	}

// 	router.Any("/v1/*any", gin.WrapH(grpcMux))

// 	router.GET("/ping", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{"message": "pong"})
// 	})

// 	server.rServer = &http.Server{
// 		// Addr:         ":9091",
// 		Addr:         server.config.GRPCGWServerAddress,
// 		Handler:      router,
// 		ReadTimeout:  R_TIME_OUT,
// 		WriteTimeout: W_TIME_OUT,
// 	}

// 	return server, nil
// }

func (server *Server) StartgPRC() error {
	logger.Log.Print(2, "gRPC server start.%s", server.config.GRPCServerAddress)

	// listener, err := net.Listen("tcp", ":50051")
	listener, err := net.Listen("tcp", server.config.GRPCServerAddress)
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

func (server *Server) ShutdowngRPC() error {
	defer server.wg.Done()
	done := make(chan struct{})
	go func() {
		server.gServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		logger.Log.Print(2, "gPRC server stopped gracefully")
	case <-time.After(5 * time.Second):
		logger.Log.Print(2, "gPRC server stopping.. timeout.. force stop")
		server.gServer.Stop()
	}

	return nil
}

func (server *Server) StartgRPCGateway() error {
	logger.Log.Print(2, "gRPC GW server start.")

	if server.rServer == nil {
		logger.Log.Error("server.rServer is null..")
	}

	if err := server.rServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Log.Error("listen error. %v", err)
		return err
	}

	return nil
}

func (server *Server) ShutdowngRPCGw() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer server.wg.Done()
	if err := server.rServer.Shutdown(ctx); err != nil {
		logger.Log.Error("Server Shutdown:", err)
		return err
	}
	return nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	addresses := strings.Split(server.config.AllowOrigins, ",")

	// router.Any("/", gin.WrapH(grpcMux))
	router.GET("/test", server.testapi)

	router.Use(corsMiddleware(addresses))
	router.Use(authMiddleware(server.tokenMaker))

	server.router = router
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
